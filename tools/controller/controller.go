package main

import (
	"encoding/json"
	"fmt"
	"github.com/jakehl/goid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/robfig/cron"
	"mini-kubernetes/tools/controller/controller_utils"
	"mini-kubernetes/tools/def"
	"mini-kubernetes/tools/etcd"
	"mini-kubernetes/tools/util"
	"os"
)

var controllerMeta = def.ControllerMeta{
	ParsedDeployments:  []*def.ParsedDeployment{},
	DeploymentNameList: []string{},
}

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	etcdClient, err := etcd.Start("", def.EtcdPort)
	controllerMeta.EtcdClient = etcdClient
	if err != nil {
		e.Logger.Error("Start etcd error!")
		os.Exit(0)
	}
	ControllerMetaInit()
	go EtcdDeploymentWatcher()
	go ReplicaChecker()
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", def.ControllerPort)))
}

func ControllerMetaInit() {
	deploymentList := controller_utils.GetDeploymentNameList(controllerMeta.EtcdClient)
	for _, name := range deploymentList {
		controllerMeta.ParsedDeployments = append(controllerMeta.ParsedDeployments, controller_utils.GetDeploymentByName(controllerMeta.EtcdClient, name))
	}
	controllerMeta.DeploymentNameList = deploymentList
}

func HandleDeploymentListChange(deploymentList []string) {
	controllerMeta.Lock.Lock()
	defer controllerMeta.Lock.Unlock()

	added, deleted := util.DifferTwoStringList(controllerMeta.DeploymentNameList, deploymentList)
	for _, name := range added {
		deployment := controller_utils.GetDeploymentByName(controllerMeta.EtcdClient, name)
		controllerMeta.ParsedDeployments = append(controllerMeta.ParsedDeployments, deployment)
		controller_utils.NewReplicaNameListByPodName(controllerMeta.EtcdClient, deployment.PodName)
		NewPodInstance(deployment.PodName, deployment.ReplicasNum)
	}
	for _, name := range deleted {
		DeleteADeployment(name)
	}
	controllerMeta.DeploymentNameList = deploymentList
}

func DeleteADeployment(name string) {
	controller_utils.RemoveAllReplicasOfPod(controllerMeta.EtcdClient, name)
	// sync cache
	for index, deployment := range controllerMeta.ParsedDeployments {
		if deployment.Name == name {
			controllerMeta.ParsedDeployments = append(controllerMeta.ParsedDeployments[:index], controllerMeta.ParsedDeployments[index+1:]...)
			break
		}
	}
}

func EtcdDeploymentWatcher() {
	watch := etcd.Watch(controllerMeta.EtcdClient, def.DeploymentListName)
	for wc := range watch {
		for _, w := range wc.Events {
			var nameList []string
			_ = json.Unmarshal(w.Kv.Value, &nameList)
			HandleDeploymentListChange(nameList)
		}
	}
}
func NewPodInstance(podName string, num int) {
	pod := controller_utils.GetPodByName(controllerMeta.EtcdClient, podName)
	for i := 0; i < num; i++ {
		podInstance := def.PodInstance{
			Pod: *pod,
			ID:  goid.NewV4UUID().String(),
			// TODO: 分配IP
			//IP:
			NodeID:        def.NodeUndefined,
			Status:        def.PENDING,
			ContainerSpec: make([]def.ContainerStatus, len(pod.Spec.Containers)),
			RestartCount:  0,
		}
		controller_utils.AddPodInstance(controllerMeta.EtcdClient, &podInstance)
	}
}

func ReplicaChecker() {
	cron2 := cron.New()
	err := cron2.AddFunc("*/10 * * * * *", CheckAllReplicas)
	if err != nil {
		fmt.Println(err)
	}
	cron2.Start()
	defer cron2.Stop()
	for {
		if controllerMeta.ShouldStop {
			break
		}
	}
}

func CheckAllReplicas() {
	controllerMeta.Lock.Lock()
	defer controllerMeta.Lock.Unlock()

	for _, deployment := range controllerMeta.ParsedDeployments {
		pod := controller_utils.GetPodByName(controllerMeta.EtcdClient, deployment.PodName)
		instancelist := controller_utils.GetReplicaNameListByPodName(controllerMeta.EtcdClient, pod.Metadata.Name)
		health := 0
		for _, instanceID := range instancelist {
			podInstance := util.GetPodInstance(instanceID, controllerMeta.EtcdClient)
			if podInstance.Status != def.FAILED {
				health++
			} else {
				controller_utils.RemovePodInstance(controllerMeta.EtcdClient, &podInstance)
			}
		}
		if health < deployment.ReplicasNum {
			NewPodInstance(pod.Metadata.Name, deployment.ReplicasNum-health)
		}
	}
}
