package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/robfig/cron"
	"mini-kubernetes/tools/def"
	"mini-kubernetes/tools/docker"
	"mini-kubernetes/tools/etcd"
	"mini-kubernetes/tools/httpget"
	"mini-kubernetes/tools/kubelet/kubelet_routines"
	net_utils "mini-kubernetes/tools/net-utils"
	//"mini-kubernetes/tools/pod"
	"mini-kubernetes/tools/resource"
	"mini-kubernetes/tools/util"
	"os"
	"strconv"
	"time"
)

var node = def.Node{}

func main() {
	parseArgs(&node.NodeName, &node.MasterIpAndPort, &node.LocalPort)
	node.NodeIP = util.GetLocalIP()
	node.ProxyPort = def.ProxyPort
	if node.NodeIP == nil {
		fmt.Println("get local ip error")
		os.Exit(0)
	}
	err := registerToMaster(&node)
	if err != nil {
		fmt.Println("network error, cannot register to master")
		os.Exit(0)
	}
	docker.CreateNetBridge(node.CniIP.String())

	/*
		creat echo instance
	*/
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	etcdClient, err := etcd.Start("", def.EtcdPort)
	if err != nil {
		e.Logger.Error("Start etcd error!")
		os.Exit(0)
	}
	node.EtcdClient = etcdClient
	cadvisorClient, err := resource.StartCadvisor()
	if err != nil {
		e.Logger.Error("Start cadvisor error!")
		os.Exit(0)
	}
	node.CadvisorClient = cadvisorClient

	//Create initial VxLANs
	net_utils.InitVxLAN(&node)

	go kubelet_routines.EtcdWatcher(&node)
	go kubelet_routines.NodesWatch(&node)
	go ResourceMonitoring()
	go ContainerCheck()

	e.Logger.Fatal(e.Start(":" + strconv.Itoa(node.LocalPort)))

}

/*
	command format:./kubelet -name `nodeName` -master `masterIP:port` -port `localPort`
	for example: ./kubelet -name node1 -master 192.168.55.184:80 -port 80
*/
func parseArgs(nodeName *string, masterIPAndPort *string, localPort *int) {
	flag.StringVar(nodeName, "name", "undefined", "name of the node, `node+nodeIP` by default")
	flag.StringVar(masterIPAndPort, "master", "undefined", "name of the node, `node+nodeIP` by default")
	flag.IntVar(localPort, "port", 100, "local port to communicate with master")
	flag.Parse()
	/*
		TODO: Check ClusterIP format legality
	*/
	if *masterIPAndPort == "undefined" {
		fmt.Println("Master Ip And Port Error!")
		os.Exit(0)
	}
}

/*
	register node to master, using http post
*/
func registerToMaster(node *def.Node) error {
	response := def.RegisterToMasterResponse{}
	request := def.RegisterToMasterRequest{
		NodeName:  node.NodeName,
		LocalIP:   node.NodeIP,
		LocalPort: node.LocalPort,
		ProxyPort: node.ProxyPort,
	}

	body, _ := json.Marshal(request)
	err, _ := httpget.Post("http://" + node.MasterIpAndPort + "/register_node").
		ContentType("application/json").
		Body(bytes.NewReader(body)).
		GetJson(&response).
		Execute()
	if err != nil {
		fmt.Println(err)
		return err
	}
	node.NodeID = response.NodeID
	node.NodeName = response.NodeName
	node.CniIP = response.CniIP

	// 为创建vxlan隧道做准备
	net_utils.InitOVS()
	return nil
}

func ContainerCheck() {
	cron2 := cron.New()
	err := cron2.AddFunc("*/5 * * * * *", checkPodRunning)
	if err != nil {
		fmt.Println(err)
	}

	cron2.Start()
	defer cron2.Stop()
	for {
		if node.ShouldStop {
			break
		}
	}
}

func IsStrInList(str string, list []string) bool {
	for _, str_in := range list {
		if str_in == str {
			return true
		}
	}
	return false
}

func checkPodRunning() {
	infos := resource.GetAllContainersInfo(node.CadvisorClient)

	var runningContainerIDs []string
	fmt.Println("infos:  ", infos)
	for _, info := range infos {
		runningContainerIDs = append(runningContainerIDs, info.Id)
	}
	fmt.Println("running container ids: ", runningContainerIDs)
	for _, instance := range node.PodInstances {
		for _, container := range instance.ContainerSpec {
			if !IsStrInList(container.ID, runningContainerIDs) {
				instance.Status = def.FAILED
				//pod.StopAndRemovePod(instance, &node)
				fmt.Println(container.ID, "fail")
				util.PersistPodInstance(*instance, node.EtcdClient)
				continue
			}
		}
	}
}

func ResourceMonitoring() {
	cron2 := cron.New()
	err := cron2.AddFunc("*/30 * * * * *", recordResource)
	if err != nil {
		fmt.Println(err)
	}

	cron2.Start()
	defer cron2.Stop()
	for {
		if node.ShouldStop {
			break
		}
	}
}

// cadvisor
func recordResource() {
	for _, podInstance := range node.PodInstances {
		if podInstance.Status != def.RUNNING {
			continue
		}
		memoryUsage := uint64(0)
		cpuLoadAverage := int32(0)
		for _, container := range podInstance.ContainerSpec {
			fmt.Println("container.ID is " + container.ID)
			info := resource.GetContainerInfoByID(node.CadvisorClient, container.ID)
			fmt.Println(info)
			memoryUsage += info.Stats[len(info.Stats)-1].Memory.Usage
			cpuLoadAverage += info.Stats[len(info.Stats)-1].Cpu.LoadAverage
		}
		key := def.GetKeyOfResourceUsageByPodInstanceID(podInstance.ID)
		resourceUsage := def.ResourceUsage{
			CPULoad:     cpuLoadAverage,
			MemoryUsage: memoryUsage,
			Time:        time.Now(),
		}
		byts, _ := json.Marshal(resourceUsage)
		etcd.Put(node.EtcdClient, key, string(byts))
	}
	key := def.KeyNodeResourceUsage(node.NodeID)
	nodeResource := resource.GetNodeResourceInfo()
	resourceUsage := def.ResourceUsage{
		CPULoad:     int32(nodeResource.TotalCPUPercent * 1000),
		MemoryUsage: nodeResource.MemoryInfo.Used,
		MemoryTotal: nodeResource.MemoryInfo.Total,
		Time:        time.Now(),
		CPUNum:      len(nodeResource.PerCPUPercent),
	}
	byts, _ := json.Marshal(resourceUsage)
	etcd.Put(node.EtcdClient, key, string(byts))
}
