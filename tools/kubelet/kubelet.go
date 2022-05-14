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
	"mini-kubernetes/tools/kubeproxy"
	net_utils "mini-kubernetes/tools/net-utils"
	"mini-kubernetes/tools/resource"
	"mini-kubernetes/tools/util"
	"net"
	"os"
	"sort"
	"time"
)

var node = def.Node{}

func main() {
	parseArgs(&node.NodeName, &node.MasterIpAndPort, &node.LocalPort)
	node.NodeIP = getLocalIP()
	node.ProxyPort = kubeproxy.ProxyPort
	if node.NodeIP == nil {
		fmt.Println("get local ip error")
		os.Exit(0)
	}
	err := registerToMaster(&node)
	if err == nil {
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

	go kubelet_routines.EtcdWatcher(&node)
	go ResourceMonitoring()
	go ContainerCheck()

	e.Logger.Fatal(e.Start(":80"))

}

/*
	command format:./kubelet --name `nodeName` --master `masterIP:port` --port `localPort`
	for example: ./kubelet --name node1 --master 192.168.55.184:80 --port 80
*/
func parseArgs(nodeName *string, masterIPAndPort *string, localPort *int) {
	flag.StringVar(nodeName, "--name", "undefined", "name of the node, `node+nodeIP` by default")
	flag.StringVar(masterIPAndPort, "--master", "undefined", "name of the node, `node+nodeIP` by default")
	flag.IntVar(localPort, "--port", 80, "local port to communicate with master")
	flag.Parse()
	/*
		TODO: Check IP format legality
	*/
	if *masterIPAndPort == "undefined" {
		fmt.Println("Master Ip And Port Error!")
		os.Exit(0)
	}
}

/*
	get local Ip
*/
func getLocalIP() net.IP {
	adds, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println(err)
		fmt.Println("cannot get local ip address, exit")
		os.Exit(0)
	}
	for _, address := range adds {
		if ip, flag_ := address.(*net.IPNet); flag_ && !ip.IP.IsLoopback() {
			if ip.IP.To4() != nil {
				return ip.IP.To4()
			}
		}
	}
	os.Exit(0)
	return nil
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

func checkPodRunning() {
	infos := resource.GetAllContainersInfo(node.CadvisorClient)
	var runningContainerIDs []string
	for _, info := range infos {
		runningContainerIDs = append(runningContainerIDs, info.Id)
	}
	sort.Strings(runningContainerIDs)
	for _, instance := range node.PodInstances {
		for _, container := range instance.ContainerSpec {
			if sort.SearchStrings(runningContainerIDs, container.ID) == len(runningContainerIDs) {
				instance.Status = def.FAILED
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

func recordResource() {
	for _, podInstance := range node.PodInstances {
		if podInstance.Status != def.RUNNING {
			continue
		}
		memoryUsage := uint64(0)
		cpuLoadAverage := int32(0)
		for _, container := range podInstance.ContainerSpec {
			info := resource.GetContainerInfoByName(node.CadvisorClient, container.Name)
			memoryUsage += info.Stats[len(info.Stats)-1].Memory.Usage
			cpuLoadAverage += info.Stats[len(info.Stats)-1].Cpu.LoadAverage
		}
		key := fmt.Sprintf("%s_resource_usage", podInstance.ID)
		resp := etcd.Get(node.EtcdClient, key)
		resourceSeq := def.ResourceUsageSequence{}
		jsonString := ``
		for _, ev := range resp.Kvs {
			jsonString += fmt.Sprintf(`"%s":"%s", `, ev.Key, ev.Value)
		}
		jsonString = fmt.Sprintf(`{%s}`, jsonString)
		_ = json.Unmarshal([]byte(jsonString), &resourceSeq)
		resourceUsage := def.ResourceUsage{
			CPULoad:     cpuLoadAverage,
			MemoryUsage: memoryUsage,
			Time:        time.Now(),
		}
		if len(resourceSeq.Sequence) < 30 {
			resourceSeq.Sequence = append(resourceSeq.Sequence, resourceUsage)
		} else {
			resourceSeq.Sequence = append(resourceSeq.Sequence[1:], resourceUsage)
		}
		byts, _ := json.Marshal(resourceSeq)
		etcd.Put(node.EtcdClient, key, string(byts))
	}
	key := fmt.Sprintf("%d_resource_usage", node.NodeID)
	resp := etcd.Get(node.EtcdClient, key)
	resourceSeq := def.ResourceUsageSequence{}
	jsonString := ``
	for _, ev := range resp.Kvs {
		jsonString += fmt.Sprintf(`"%s":"%s"`, ev.Key, ev.Value)
	}
	jsonString = fmt.Sprintf(`{%s}`, jsonString)
	_ = json.Unmarshal([]byte(jsonString), &resourceSeq)
	nodeResource := resource.GetNodeResourceInfo()
	resourceUsage := def.ResourceUsage{
		CPULoad:     int32(nodeResource.TotalCPUPercent * 1000),
		MemoryUsage: nodeResource.MemoryInfo.Used,
		MemoryTotal: nodeResource.MemoryInfo.Total,
		Time:        time.Now(),
	}
	if len(resourceSeq.Sequence) < 30 {
		resourceSeq.Sequence = append(resourceSeq.Sequence, resourceUsage)
	} else {
		resourceSeq.Sequence = append(resourceSeq.Sequence[1:], resourceUsage)
	}
	byts, _ := json.Marshal(resourceSeq)
	etcd.Put(node.EtcdClient, key, string(byts))
}
