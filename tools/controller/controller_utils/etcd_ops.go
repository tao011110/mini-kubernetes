package controller_utils

import (
	"encoding/json"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"mini-kubernetes/tools/def"
	"mini-kubernetes/tools/etcd"
	"mini-kubernetes/tools/util"
)

func GetDeploymentNameList(etcdClient *clientv3.Client) []string {
	var deploymentNameList []string
	util.EtcdUnmarshal(etcd.Get(etcdClient, def.DeploymentListName), &deploymentNameList)
	return deploymentNameList
}

func GetHorizontalPodAutoscalerNameList(etcdClient *clientv3.Client) []string {
	var horizontalPodAutoscalerNameList []string
	util.EtcdUnmarshal(etcd.Get(etcdClient, def.HorizontalPodAutoscalerListName), &horizontalPodAutoscalerNameList)
	return horizontalPodAutoscalerNameList
}

func GetDeploymentByName(etcdClient *clientv3.Client, deploymentName string) *def.ParsedDeployment {
	key := def.GetKeyOfDeployment(deploymentName)
	deployment := def.ParsedDeployment{}
	util.EtcdUnmarshal(etcd.Get(etcdClient, key), &deployment)
	return &deployment
}

func GetHorizontalPodAutoscalerByName(etcdClient *clientv3.Client, horizontalPodAutoscalerName string) *def.ParsedHorizontalPodAutoscaler {
	key := def.GetKeyOfAutoscaler(horizontalPodAutoscalerName)
	horizontalPodAutoscaler := def.ParsedHorizontalPodAutoscaler{}
	util.EtcdUnmarshal(etcd.Get(etcdClient, key), &horizontalPodAutoscaler)
	return &horizontalPodAutoscaler
}

func AddPodInstance(etcdClient *clientv3.Client, instance *def.PodInstance) {
	util.PersistPodInstance(*instance, etcdClient)
	{
		key := def.GetKeyOfPodReplicasNameListByPodName(instance.Pod.Metadata.Name)
		var instanceIDList []string
		util.EtcdUnmarshal(etcd.Get(etcdClient, key), &instanceIDList)
		instanceIDList = append(instanceIDList, instance.ID)
		newJsonString, _ := json.Marshal(instanceIDList)
		etcd.Put(etcdClient, key, string(newJsonString))
	}
	{
		var instanceIDList []string
		util.EtcdUnmarshal(etcd.Get(etcdClient, def.PodInstanceListID), &instanceIDList)
		instanceIDList = append(instanceIDList, instance.ID)
		newJsonString, _ := json.Marshal(instanceIDList)
		etcd.Put(etcdClient, def.PodInstanceListID, string(newJsonString))
	}
	{
		resourceUsage := def.ResourceUsage{
			Valid: false,
		}
		newJsonString, _ := json.Marshal(resourceUsage)
		etcd.Put(etcdClient, def.GetKeyOfResourceUsageByPodInstanceID(instance.ID), string(newJsonString))
	}
}

func RemovePodInstance(etcdClient *clientv3.Client, instance *def.PodInstance) {
	{
		var podInstanceIDList []string
		util.EtcdUnmarshal(etcd.Get(etcdClient, def.PodInstanceListID), &podInstanceIDList)
		for index, id := range podInstanceIDList {
			if id == instance.ID {
				podInstanceIDList = append(podInstanceIDList[:index], podInstanceIDList[index+1:]...)
				break
			}
		}
		fmt.Println("instance.ID:  ", instance.ID)
		fmt.Println("after delete, podInstanceIDList:  ", podInstanceIDList)
		newJsonString, _ := json.Marshal(podInstanceIDList)
		etcd.Put(etcdClient, def.PodInstanceListID, string(newJsonString))
	}
	{
		key := def.GetKeyOfPodReplicasNameListByPodName(instance.Pod.Metadata.Name)
		var podInstanceNameList []string
		util.EtcdUnmarshal(etcd.Get(etcdClient, key), &podInstanceNameList)
		for index, id := range podInstanceNameList {
			if id == instance.ID {
				podInstanceNameList = append(podInstanceNameList[:index], podInstanceNameList[index+1:]...)
				break
			}
		}
		newJsonString, _ := json.Marshal(podInstanceNameList)
		etcd.Put(etcdClient, key, string(newJsonString))
	}
}

func RemoveAllReplicasOfPod(etcdClient *clientv3.Client, podName string) {
	// remove from instance list, scheduler will remove it from node
	fmt.Println("podName is:  ", podName)
	key := def.GetKeyOfPodReplicasNameListByPodName(podName)
	podInstanceIDList := GetReplicaNameListByPodName(etcdClient, podName)
	fmt.Println("GetReplicaNameListByPodName: ", key)
	fmt.Println("podInstanceIDList: ", podInstanceIDList)
	for _, instanceID := range podInstanceIDList {
		fmt.Println("try to get by instanceID:  ", instanceID)
		instance := util.GetPodInstance(instanceID, etcdClient)
		RemovePodInstance(etcdClient, &instance)
	}
	// remove it's pod-replica entry
	etcd.Delete(etcdClient, key)
}

func GetReplicaNameListByPodName(etcdClient *clientv3.Client, podName string) []string {
	key := def.GetKeyOfPodReplicasNameListByPodName(podName)
	var deploymentNameList []string
	util.EtcdUnmarshal(etcd.Get(etcdClient, key), &deploymentNameList)
	return deploymentNameList
}

func NewReplicaNameListByPodName(etcdClient *clientv3.Client, podName string) {
	// add empty list
	key := def.GetKeyOfPodReplicasNameListByPodName(podName)
	var deploymentNameList []string
	newJsonString, _ := json.Marshal(deploymentNameList)
	etcd.Put(etcdClient, key, string(newJsonString))
}

func GetPodByName(etcdClient *clientv3.Client, podName string) *def.Pod {
	pod := def.Pod{}
	key := def.GetKeyOfPod(podName)
	util.EtcdUnmarshal(etcd.Get(etcdClient, key), &pod)
	return &pod
}

func GetPodInstanceResourceUsageByName(etcdClient *clientv3.Client, podInstanceID string) *def.ResourceUsage {
	key := def.GetKeyOfResourceUsageByPodInstanceID(podInstanceID)
	resourceUsage := def.ResourceUsage{}
	util.EtcdUnmarshal(etcd.Get(etcdClient, key), &resourceUsage)
	return &resourceUsage
}
