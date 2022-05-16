package controller_utils

import (
	"encoding/json"
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
	deployment := def.ParsedDeployment{}
	util.EtcdUnmarshal(etcd.Get(etcdClient, deploymentName), &deployment)
	return &deployment
}

func GetHorizontalPodAutoscalerByName(etcdClient *clientv3.Client, horizontalPodAutoscalerName string) *def.ParsedHorizontalPodAutoscaler {
	horizontalPodAutoscaler := def.ParsedHorizontalPodAutoscaler{}
	util.EtcdUnmarshal(etcd.Get(etcdClient, horizontalPodAutoscalerName), &horizontalPodAutoscaler)
	return &horizontalPodAutoscaler
}

func AddPodInstance(etcdClient *clientv3.Client, instance *def.PodInstance) {
	util.PersistPodInstance(*instance, etcdClient)
	{
		key := def.GetKeyOfPodReplicasNameListByPodName(instance.Pod.Metadata.Name)
		var instanceNameList []string
		util.EtcdUnmarshal(etcd.Get(etcdClient, key), &instanceNameList)
		instanceNameList = append(instanceNameList, instance.ID)
		newJsonString, _ := json.Marshal(instanceNameList)
		etcd.Put(etcdClient, key, string(newJsonString))
	}
	{
		var instanceNameList []string
		util.EtcdUnmarshal(etcd.Get(etcdClient, def.PodInstanceListName), &instanceNameList)
		instanceNameList = append(instanceNameList, instance.ID)
		newJsonString, _ := json.Marshal(instanceNameList)
		etcd.Put(etcdClient, def.PodInstanceListName, string(newJsonString))
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
		var deploymentNameList []string
		util.EtcdUnmarshal(etcd.Get(etcdClient, def.PodInstanceListName), &deploymentNameList)
		for index, name := range deploymentNameList {
			if name == instance.ID {
				deploymentNameList = append(deploymentNameList[:index], deploymentNameList[index+1:]...)
				break
			}
		}
		newJsonString, _ := json.Marshal(deploymentNameList)
		etcd.Put(etcdClient, def.PodInstanceListName, string(newJsonString))
	}
	{
		key := def.GetKeyOfPodReplicasNameListByPodName(instance.Pod.Metadata.Name)
		var deploymentNameList []string
		util.EtcdUnmarshal(etcd.Get(etcdClient, key), &deploymentNameList)
		for index, name := range deploymentNameList {
			if name == instance.ID {
				deploymentNameList = append(deploymentNameList[:index], deploymentNameList[index+1:]...)
				break
			}
		}
		newJsonString, _ := json.Marshal(deploymentNameList)
		etcd.Put(etcdClient, key, string(newJsonString))
	}
}

func RemoveAllReplicasOfPod(etcdClient *clientv3.Client, podName string) {
	// remove from instance list, scheduler will remove it from node
	key := def.GetKeyOfPodReplicasNameListByPodName(podName)
	deploymentNameList := GetReplicaNameListByPodName(etcdClient, podName)
	for _, instanceName := range deploymentNameList {
		instance := util.GetPodInstance(instanceName, etcdClient)
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
	util.EtcdUnmarshal(etcd.Get(etcdClient, podName), &pod)
	return &pod
}

func GetPodInstanceResourceUsageByName(etcdClient *clientv3.Client, podInstanceID string) *def.ResourceUsage {
	key := def.GetKeyOfResourceUsageByPodInstanceID(podInstanceID)
	resourceUsage := def.ResourceUsage{}
	util.EtcdUnmarshal(etcd.Get(etcdClient, key), &resourceUsage)
	return &resourceUsage
}
