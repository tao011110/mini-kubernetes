package main

import (
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/robfig/cron"
	"math"
	"mini-kubernetes/tools/controller/controller_utils"
	"mini-kubernetes/tools/def"
	"mini-kubernetes/tools/etcd"
	"mini-kubernetes/tools/util"
	"os"
)

var controllerMeta = def.ControllerMeta{
	ParsedDeployments:                []*def.ParsedDeployment{},
	DeploymentNameList:               []string{},
	ParsedHorizontalPodAutoscalers:   []*def.ParsedHorizontalPodAutoscaler{},
	HorizontalPodAutoscalersNameList: []string{},
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
	go EtcdHorizontalPodAutoscalerWatcher()

	go ReplicaChecker()
	go HorizontalPodAutoscalerChecker()
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", def.ControllerPort)))
}

func ControllerMetaInit() {
	deploymentList := controller_utils.GetDeploymentNameList(controllerMeta.EtcdClient)
	horizontalPodAutoscalerNameList := controller_utils.GetHorizontalPodAutoscalerNameList(controllerMeta.EtcdClient)
	for _, name := range deploymentList {
		controllerMeta.ParsedDeployments = append(controllerMeta.ParsedDeployments, controller_utils.GetDeploymentByName(controllerMeta.EtcdClient, name))
	}
	for _, name := range horizontalPodAutoscalerNameList {
		controllerMeta.ParsedHorizontalPodAutoscalers = append(controllerMeta.ParsedHorizontalPodAutoscalers, controller_utils.GetHorizontalPodAutoscalerByName(controllerMeta.EtcdClient, name))
	}
	controllerMeta.DeploymentNameList = deploymentList
	controllerMeta.HorizontalPodAutoscalersNameList = horizontalPodAutoscalerNameList
}

func HandleDeploymentListChange(deploymentList []string) {
	controllerMeta.DeploymentLock.Lock()
	defer controllerMeta.DeploymentLock.Unlock()

	fmt.Println(controllerMeta.DeploymentNameList)
	fmt.Println(deploymentList)
	added, deleted := util.DifferTwoStringList(controllerMeta.DeploymentNameList, deploymentList)
	fmt.Println("added:   ", added)
	fmt.Println("deleted:   ", deleted)
	for _, name := range added {
		deployment := controller_utils.GetDeploymentByName(controllerMeta.EtcdClient, name)
		fmt.Println(deployment)
		controllerMeta.ParsedDeployments = append(controllerMeta.ParsedDeployments, deployment)
		controller_utils.NewReplicaNameListByPodName(controllerMeta.EtcdClient, deployment.PodName)
		controller_utils.NewNPodInstance(controllerMeta.EtcdClient, deployment.PodName, deployment.ReplicasNum)
	}
	for _, name := range deleted {
		DeleteADeployment(name)
	}
	controllerMeta.DeploymentNameList = deploymentList
}

func HandleHorizontalPodAutoscalerListChange(horizontalPodAutoscalerList []string) {
	controllerMeta.HorizontalPodAutoscalersLock.Lock()
	defer controllerMeta.HorizontalPodAutoscalersLock.Unlock()

	added, deleted := util.DifferTwoStringList(controllerMeta.DeploymentNameList, horizontalPodAutoscalerList)
	for _, name := range added {
		horizontalPodAutoscaler := controller_utils.GetHorizontalPodAutoscalerByName(controllerMeta.EtcdClient, name)
		controllerMeta.ParsedHorizontalPodAutoscalers = append(controllerMeta.ParsedHorizontalPodAutoscalers, horizontalPodAutoscaler)
		controller_utils.NewReplicaNameListByPodName(controllerMeta.EtcdClient, horizontalPodAutoscaler.PodName)
		controller_utils.NewNPodInstance(controllerMeta.EtcdClient, horizontalPodAutoscaler.PodName, horizontalPodAutoscaler.MinReplicas)
	}
	for _, name := range deleted {
		DeleteAHorizontalPodAutoscaler(name)
	}
	controllerMeta.HorizontalPodAutoscalersNameList = horizontalPodAutoscalerList
}

func DeleteAHorizontalPodAutoscaler(name string) {
	controller_utils.RemoveAllReplicasOfPod(controllerMeta.EtcdClient, def.GetPodNameOfAutoscaler(name))
	// sync cache
	for index, horizontalPodAutoscaler := range controllerMeta.ParsedHorizontalPodAutoscalers {
		if horizontalPodAutoscaler.Name == name {
			controllerMeta.ParsedHorizontalPodAutoscalers = append(controllerMeta.ParsedHorizontalPodAutoscalers[:index], controllerMeta.ParsedHorizontalPodAutoscalers[index+1:]...)
			break
		}
	}
}

func DeleteADeployment(name string) {
	controller_utils.RemoveAllReplicasOfPod(controllerMeta.EtcdClient, def.GetPodNameOfDeployment(name))
	// sync cache
	for index, deployment := range controllerMeta.ParsedDeployments {
		if deployment.Name == name {
			controllerMeta.ParsedDeployments = append(controllerMeta.ParsedDeployments[:index], controllerMeta.ParsedDeployments[index+1:]...)
			break
		}
	}
}

func EtcdHorizontalPodAutoscalerWatcher() {
	watch := etcd.Watch(controllerMeta.EtcdClient, def.HorizontalPodAutoscalerListName)
	for wc := range watch {
		for _, w := range wc.Events {
			var nameList []string
			_ = json.Unmarshal(w.Kv.Value, &nameList)
			HandleHorizontalPodAutoscalerListChange(nameList)
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

func ReplicaChecker() {
	cron2 := cron.New()
	err := cron2.AddFunc("*/5 * * * * *", CheckAllReplicas)
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

func HorizontalPodAutoscalerChecker() {
	cron2 := cron.New()
	err := cron2.AddFunc("*/30 * * * * *", CheckAllHorizontalPodAutoscalers)
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
	controllerMeta.DeploymentLock.Lock()
	defer controllerMeta.DeploymentLock.Unlock()

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
			controller_utils.NewNPodInstance(controllerMeta.EtcdClient, pod.Metadata.Name, deployment.ReplicasNum-health)
		}
	}
}

func CheckAllHorizontalPodAutoscalers() {
	controllerMeta.HorizontalPodAutoscalersLock.Lock()
	defer controllerMeta.HorizontalPodAutoscalersLock.Unlock()

	for _, horizontalPodAutoscaler := range controllerMeta.ParsedHorizontalPodAutoscalers {
		pod := controller_utils.GetPodByName(controllerMeta.EtcdClient, horizontalPodAutoscaler.PodName)
		instancelist := controller_utils.GetReplicaNameListByPodName(controllerMeta.EtcdClient, pod.Metadata.Name)
		cpu := float64(0)
		memory := int64(0)
		minCPUUsagePodInstance := def.PodInstance{}
		minCPUUsage := math.MaxFloat64
		minMemoryUsagePodInstance := def.PodInstance{}
		minMemoryUsage := int64(math.MaxInt64)
		activeNum := 0
		for _, instanceID := range instancelist {
			podInstance := util.GetPodInstance(instanceID, controllerMeta.EtcdClient)
			if podInstance.Status == def.FAILED {
				continue
			} else {
				activeNum++
				podInstanceResourceUsage := controller_utils.GetPodInstanceResourceUsageByName(controllerMeta.EtcdClient, instanceID)
				if podInstanceResourceUsage.Valid {
					instanceCPUUsage := float64(podInstanceResourceUsage.CPULoad) / 1000
					cpu += instanceCPUUsage
					if instanceCPUUsage < minCPUUsage {
						minCPUUsage = instanceCPUUsage
						minCPUUsagePodInstance = podInstance
					}

					instanceMemoryUsage := int64(podInstanceResourceUsage.MemoryUsage)
					memory += instanceMemoryUsage
					if instanceMemoryUsage < minMemoryUsage {
						minMemoryUsage = instanceMemoryUsage
						minMemoryUsagePodInstance = podInstance
					}
				}
			}
		}
		// TODO: CPU字段含义有点问题
		// TODO: 优化调度策略
		// 优先保证最低需求
		if activeNum < horizontalPodAutoscaler.MinReplicas {
			controller_utils.NewNPodInstance(controllerMeta.EtcdClient, horizontalPodAutoscaler.PodName, horizontalPodAutoscaler.MinReplicas-activeNum)
		} else if cpu < 0.8*horizontalPodAutoscaler.CPUMinValue || float64(memory) < 0.8*float64(horizontalPodAutoscaler.MemoryMinValue) {
			if activeNum < horizontalPodAutoscaler.MaxReplicas {
				controller_utils.NewNPodInstance(controllerMeta.EtcdClient, horizontalPodAutoscaler.PodName, 1)
			}
		} else if cpu > 1.2*horizontalPodAutoscaler.CPUMinValue {
			if activeNum > horizontalPodAutoscaler.MinReplicas {
				controller_utils.RemovePodInstance(controllerMeta.EtcdClient, &minCPUUsagePodInstance)
			}
		} else if float64(memory) > 1.2*float64(horizontalPodAutoscaler.MemoryMinValue) {
			if activeNum > horizontalPodAutoscaler.MinReplicas {
				controller_utils.RemovePodInstance(controllerMeta.EtcdClient, &minMemoryUsagePodInstance)
			}
		}
	}
}
