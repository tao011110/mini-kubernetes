package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/monaco-io/request"
	"github.com/robfig/cron"
	"mini-kubernetes/tools/def"
	"mini-kubernetes/tools/docker"
	"mini-kubernetes/tools/etcd"
	"mini-kubernetes/tools/httpget"
	"mini-kubernetes/tools/kubelet/kubelet_routines"
	"mini-kubernetes/tools/kubelet/kubelet_utils"
	net_utils "mini-kubernetes/tools/net-utils"
	"mini-kubernetes/tools/pod"
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
	KubeletInitialize()
	go EtcdWatcher()
	go kubelet_routines.NodesWatch(node.NodeID, node.EtcdClient)
	go ResourceMonitoring()
	go ContainerCheck()
	go HeartBeatSender()

	e.Logger.Fatal(e.Start(":" + strconv.Itoa(node.LocalPort)))

}

func KubeletInitialize() {
	currentPodInstanceList := kubelet_utils.GetAllPodInstancesOfANode(node.NodeID, node.EtcdClient)
	for _, instanceID := range currentPodInstanceList {
		instance := util.GetPodInstance(instanceID, node.EtcdClient)
		if instance.Status == def.RUNNING || instance.Status == def.FAILED || instance.Status == def.SUCCEEDED {
			node.PodInstances = append(node.PodInstances, &instance)
		}
	}
	handlePodInstancesChange(currentPodInstanceList)
}

/*
	command format:./kubelet -name `nodeName` -master `masterIP:port` -port `localPort`
	for example: ./kubelet -name node1 -master 10.119.11.140:8000 -port 100
*/
func parseArgs(nodeName *string, masterIPAndPort *string, localPort *int) {
	flag.StringVar(nodeName, "name", "undefined", "name of the node, `node+nodeIP` by default")
	flag.StringVar(masterIPAndPort, "master", "undefined", "name of the node, `node+nodeIP` by default")
	flag.IntVar(localPort, "port", 100, "local port to communicate with master")
	flag.Parse()
	if *masterIPAndPort == "undefined" {
		fmt.Println("Master Ip And Port Error!")
		os.Exit(0)
	}
}

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

func checkPodRunning() {
	infos := resource.GetAllContainersInfo(node.CadvisorClient)
	var runningContainerIDs []string
	fmt.Println("infos:  ", infos)
	for _, info := range infos {
		runningContainerIDs = append(runningContainerIDs, info.Id)
	}
	fmt.Println("running container ids: ", runningContainerIDs)
	for _, instance := range node.PodInstances {
		if instance.Status != def.RUNNING {
			continue
		}
		for _, container := range instance.ContainerSpec {
			if !kubelet_utils.IsStrInList(container.ID, runningContainerIDs) {
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
	err := cron2.AddFunc("*/15 * * * * *", recordResource)
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
		kubelet_utils.RecordPodInstanceResource(*podInstance, node.CadvisorClient, node.EtcdClient)
	}
	kubelet_utils.RecordNodeResource(node.NodeID, node.EtcdClient)
}

func EtcdWatcher() {
	key := def.GetKeyOfPodInstanceListKeyOfNodeByID(node.NodeID)
	watch := etcd.Watch(node.EtcdClient, key)
	for wc := range watch {
		for _, w := range wc.Events {
			var instances []string
			_ = json.Unmarshal(w.Kv.Value, &instances)
			handlePodInstancesChange(instances)
		}
	}
}

func handlePodInstancesChange(instances []string) {
	var instancesCurrent []string
	for _, instance := range node.PodInstances {
		instancesCurrent = append(instancesCurrent, instance.ID)
	}
	adds, deletedIDs := util.DifferTwoStringList(instancesCurrent, instances)
	for _, delete_ := range deletedIDs {
		for index, instance := range node.PodInstances {
			if delete_ == instance.ID && instance.Status == def.RUNNING {
				pod.StopAndRemovePod(node.PodInstances[index], &node)
				break
			}
		}
	}
	for _, add := range adds {
		fmt.Println("add:   ", add)
		podInstance := util.GetPodInstance(add, node.EtcdClient)
		if podInstance.Status != def.PENDING {
			continue
		}
		pod.CreateAndStartPod(&podInstance, &node)
		node.PodInstances = append(node.PodInstances, &podInstance)
	}
}

func HeartBeatSender() {
	cron2 := cron.New()
	err := cron2.AddFunc("*/60 * * * * *", SendHeartBeat)
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

func SendHeartBeat() {
	c := request.Client{
		URL:    fmt.Sprintf("http://%s/heartbeat", node.MasterIpAndPort),
		Method: "POST",
		JSON: def.HeartBeat{
			NodeID:    node.NodeID,
			TimeStamp: time.Now(),
		},
	}
	_ = c.Send()
}
