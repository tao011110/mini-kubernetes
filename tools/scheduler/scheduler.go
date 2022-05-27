package main

import (
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/robfig/cron"
	"mini-kubernetes/tools/def"
	"mini-kubernetes/tools/etcd"
	"mini-kubernetes/tools/scheduler/scheduler_utils"
	"mini-kubernetes/tools/util"
	"os"
	"strconv"
)

var schedulerMeta = def.Scheduler{
	ScheduledPodInstancesName: []string{},
	Nodes:                     []*def.NodeInfoSchedulerCache{},
	CannotSchedule:            []string{},
	ShouldStop:                false,
}

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	etcdClient, err := etcd.Start("", def.EtcdPort)
	schedulerMeta.EtcdClient = etcdClient
	if err != nil {
		e.Logger.Error("Start etcd error!")
		os.Exit(0)
	}
	SchedulerMetaInit()
	go EtcdNodeWatcher()
	go EtcdPodInstanceWatcher()
	go ReScheduleCannotScheduleInstanceRoutine()
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", def.SchedulerPort)))
}

func SchedulerMetaInit() {
	nodeIDList := scheduler_utils.GetAllNodesID(schedulerMeta.EtcdClient)
	for _, nodeID := range nodeIDList {
		nodeInfo, replicasOnNode := scheduler_utils.GetInfoOfANode(schedulerMeta.EtcdClient, nodeID)
		schedulerMeta.Nodes = append(schedulerMeta.Nodes, nodeInfo)
		schedulerMeta.ScheduledPodInstancesName = append(schedulerMeta.ScheduledPodInstancesName, replicasOnNode...)
	}
	currentAllReplicas := scheduler_utils.GetAllPodInstancesID(schedulerMeta.EtcdClient)
	HandlePodInstanceChange(currentAllReplicas)
}

func EtcdNodeWatcher() {
	watch := etcd.Watch(schedulerMeta.EtcdClient, def.NodeListID)
	for wc := range watch {
		for _, w := range wc.Events {
			var instances []int
			_ = json.Unmarshal(w.Kv.Value, &instances)
			HandleNodeListChange(instances)
		}
	}
}

func HandleNodeListChange(newNodeList []int) {
	schedulerMeta.Lock.Lock()
	defer schedulerMeta.Lock.Unlock()

	news := util.ConvertIntListToStringList(newNodeList)
	var olds []string
	for _, old := range schedulerMeta.Nodes {
		olds = append(olds, fmt.Sprintf("%d", old.NodeID))
	}
	added, deleted := util.DifferTwoStringList(olds, news)
	fmt.Println("nodes  added:  ", added)
	fmt.Println("nodes  deleted", deleted)
	for _, add := range added {
		nodeId, _ := strconv.Atoi(add)
		nodeInfo, replicas := scheduler_utils.GetInfoOfANode(schedulerMeta.EtcdClient, nodeId)
		schedulerMeta.Nodes = append(schedulerMeta.Nodes, nodeInfo)
		schedulerMeta.ScheduledPodInstancesName = append(schedulerMeta.ScheduledPodInstancesName, replicas...)
	}
	for _, delete_ := range deleted {
		nodeId, _ := strconv.Atoi(delete_)
		for index, node := range schedulerMeta.Nodes {
			if node.NodeID == nodeId {
				schedulerMeta.Nodes = append(schedulerMeta.Nodes[:index], schedulerMeta.Nodes[index+1:]...)
				break
			}
		}
	}
}

func EtcdPodInstanceWatcher() {
	watch := etcd.Watch(schedulerMeta.EtcdClient, def.PodInstanceListID)
	for wc := range watch {
		for _, w := range wc.Events {
			var instances []string
			_ = json.Unmarshal(w.Kv.Value, &instances)
			HandlePodInstanceChange(instances)
		}
	}
}

func HandlePodInstanceChange(news []string) {
	schedulerMeta.Lock.Lock()
	defer schedulerMeta.Lock.Unlock()

	fmt.Println(schedulerMeta.ScheduledPodInstancesName)
	fmt.Println(news)
	addeds, deleteds := util.DifferTwoStringList(schedulerMeta.ScheduledPodInstancesName, news)
	fmt.Println("addeds:  ", addeds)
	fmt.Println("deleteds:  ", deleteds)
	for _, deleted := range deleteds {
		found := false
		for _, node := range schedulerMeta.Nodes {
			for _, replica := range node.PodInstanceList {
				if replica.InstanceName == deleted {
					found = true
					scheduler_utils.DeletePodInstanceFromNode(schedulerMeta.EtcdClient, node.NodeID, deleted)
					DeletePodInstanceFromSchedulerCache(node.NodeID, deleted)
					break
				}
			}
			if found {
				break
			}
		}
	}
	for _, added := range addeds {
		success, nodeID, podInstance := SchedulePodInstanceToNode(added)
		if success {
			fmt.Printf("success\n")
			scheduler_utils.AddPodInstanceToNode(schedulerMeta.EtcdClient, nodeID, podInstance)
			AddPodInstanceFromSchedulerCache(nodeID, podInstance)
			schedulerMeta.ScheduledPodInstancesName = append(schedulerMeta.ScheduledPodInstancesName, added)
		} else {
			fmt.Printf("failed\n")
			schedulerMeta.CannotSchedule = append(schedulerMeta.CannotSchedule, added)
		}
	}
}

func DeletePodInstanceFromSchedulerCache(nodeID int, instanceName string) {
	for _, nodeInfo := range schedulerMeta.Nodes {
		if nodeID == nodeInfo.NodeID {
			for index, podInstanceInfo := range nodeInfo.PodInstanceList {
				if podInstanceInfo.InstanceName == instanceName {
					nodeInfo.PodInstanceList = append(nodeInfo.PodInstanceList[:index], nodeInfo.PodInstanceList[index+1:]...)
					break
				}
			}
			break
		}
	}
}

func AddPodInstanceFromSchedulerCache(nodeID int, instance *def.PodInstance) {
	for _, nodeInfo := range schedulerMeta.Nodes {
		if nodeID == nodeInfo.NodeID {
			newInstanceCache := def.PodInstanceSchedulerCache{
				InstanceName: instance.ID,
				PodName:      instance.Pod.Metadata.Name,
			}
			nodeInfo.PodInstanceList = append(nodeInfo.PodInstanceList, newInstanceCache)
			break
		}
	}
}

func SchedulePodInstanceToNode(instanceID string) (success bool, nodeId int, podInstance *def.PodInstance) {
	{
		podInstanceTmp := util.GetPodInstance(instanceID, schedulerMeta.EtcdClient)
		podInstance = &podInstanceTmp
	}
	var nodeIndexList []int
	for index := range schedulerMeta.Nodes {
		nodeIndexList = append(nodeIndexList, index)
	}
	fmt.Println("nodeIndexList:  ", nodeIndexList)
	notWithFilterResult := scheduler_utils.NotWithFilter(nodeIndexList, podInstance.NodeSelector.NotWith, schedulerMeta.Nodes)
	if len(notWithFilterResult) == 0 {
		success = false
		return
	}
	fmt.Println("[notWithFilterResult]", notWithFilterResult)
	withFilterResult := scheduler_utils.WithFilter(notWithFilterResult, podInstance.NodeSelector.With, schedulerMeta.Nodes)
	if len(withFilterResult) == 0 {
		success = false
		return
	}
	fmt.Println("[withFilterResult]", withFilterResult)
	CPU, memory := scheduler_utils.PodResourceRequest(podInstance)
	// TODO: memory的单位
	resourceFilterResult := scheduler_utils.ResourceFilter(schedulerMeta.EtcdClient,
		withFilterResult,
		CPU, memory, schedulerMeta.Nodes)
	fmt.Println("[resourceFilterResult]", resourceFilterResult)
	if len(resourceFilterResult) == 0 {
		success = false
		return
	}
	success = true
	choseNodeIndex := scheduler_utils.ChooseNode(schedulerMeta.EtcdClient, resourceFilterResult, schedulerMeta.Nodes)
	fmt.Println("chose node index is ", choseNodeIndex, " id is ", schedulerMeta.Nodes[choseNodeIndex].NodeID)
	nodeId = schedulerMeta.Nodes[choseNodeIndex].NodeID
	fmt.Println("Schedule to node:  ", nodeId)
	return
}

func ReScheduleCannotScheduleInstanceRoutine() {
	cron2 := cron.New()
	err := cron2.AddFunc("*/15 * * * * *", ReSchedulerCannotScheduleInstance)
	if err != nil {
		fmt.Println(err)
	}
	cron2.Start()
	defer cron2.Stop()
	for {
		if schedulerMeta.ShouldStop {
			break
		}
	}
}

func ReSchedulerCannotScheduleInstance() {
	schedulerMeta.Lock.Lock()
	defer schedulerMeta.Lock.Unlock()

	var newFailReSchedule []string
	for _, instanceName := range schedulerMeta.CannotSchedule {
		success, nodeID, podInstance := SchedulePodInstanceToNode(instanceName)
		if success {
			scheduler_utils.AddPodInstanceToNode(schedulerMeta.EtcdClient, nodeID, podInstance)
			AddPodInstanceFromSchedulerCache(nodeID, podInstance)
			schedulerMeta.ScheduledPodInstancesName = append(schedulerMeta.ScheduledPodInstancesName, instanceName)
		} else {
			newFailReSchedule = append(newFailReSchedule, instanceName)
		}
	}
	schedulerMeta.CannotSchedule = newFailReSchedule
}
