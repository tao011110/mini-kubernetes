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
	resp := etcd.Get(etcdClient, def.DeploymentListName)
	var deploymentNameList []string
	jsonString := ``
	for _, ev := range resp.Kvs {
		jsonString += fmt.Sprintf(`"%s":"%s"`, ev.Key, ev.Value)
	}
	jsonString = fmt.Sprintf(`{%s}`, jsonString)
	_ = json.Unmarshal([]byte(jsonString), &deploymentNameList)
	return deploymentNameList
}

func GetHorizontalPodAutoscalerNameList(etcdClient *clientv3.Client) []string {
	resp := etcd.Get(etcdClient, def.HorizontalPodAutoscalerListName)
	var horizontalPodAutoscalerNameList []string
	jsonString := ``
	for _, ev := range resp.Kvs {
		jsonString += fmt.Sprintf(`"%s":"%s"`, ev.Key, ev.Value)
	}
	jsonString = fmt.Sprintf(`{%s}`, jsonString)
	_ = json.Unmarshal([]byte(jsonString), &horizontalPodAutoscalerNameList)
	return horizontalPodAutoscalerNameList
}

func GetDeploymentByName(etcdClient *clientv3.Client, deploymentName string) *def.ParsedDeployment {
	resp := etcd.Get(etcdClient, deploymentName)
	deployment := def.ParsedDeployment{}
	jsonString := ``
	for _, ev := range resp.Kvs {
		jsonString += fmt.Sprintf(`"%s":"%s"`, ev.Key, ev.Value)
	}
	jsonString = fmt.Sprintf(`{%s}`, jsonString)
	_ = json.Unmarshal([]byte(jsonString), &deployment)
	return &deployment
}

func GetHorizontalPodAutoscalerByName(etcdClient *clientv3.Client, horizontalPodAutoscalerName string) *def.ParsedHorizontalPodAutoscaler {
	resp := etcd.Get(etcdClient, horizontalPodAutoscalerName)
	horizontalPodAutoscaler := def.ParsedHorizontalPodAutoscaler{}
	jsonString := ``
	for _, ev := range resp.Kvs {
		jsonString += fmt.Sprintf(`"%s":"%s"`, ev.Key, ev.Value)
	}
	jsonString = fmt.Sprintf(`{%s}`, jsonString)
	_ = json.Unmarshal([]byte(jsonString), &horizontalPodAutoscaler)
	return &horizontalPodAutoscaler
}

func AddPodInstance(etcdClient *clientv3.Client, instance *def.PodInstance) {
	util.PersistPodInstance(*instance, etcdClient)
	{
		key := def.GetKeyOfPodReplicasNameListByPodName(instance.Pod.Metadata.Name)
		resp := etcd.Get(etcdClient, key) // pod -> replica_list(all replica's id of the pod)
		var instanceNameList []string
		jsonString := ``
		for _, ev := range resp.Kvs {
			jsonString += fmt.Sprintf(`"%s":"%s"`, ev.Key, ev.Value)
		}
		jsonString = fmt.Sprintf(`{%s}`, jsonString)
		_ = json.Unmarshal([]byte(jsonString), &instanceNameList)
		instanceNameList = append(instanceNameList, instance.ID)
		newJsonString, _ := json.Marshal(instanceNameList)
		etcd.Put(etcdClient, key, string(newJsonString))
	}
	{
		resp := etcd.Get(etcdClient, def.PodInstanceListName) // only 1, all podInstance's id podInstance's id : pod'name + uuid
		var instanceNameList []string
		jsonString := ``
		for _, ev := range resp.Kvs {
			jsonString += fmt.Sprintf(`"%s":"%s"`, ev.Key, ev.Value)
		}
		jsonString = fmt.Sprintf(`{%s}`, jsonString)
		_ = json.Unmarshal([]byte(jsonString), &instanceNameList)
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
		resp := etcd.Get(etcdClient, def.PodInstanceListName)
		var deploymentNameList []string
		jsonString := ``
		for _, ev := range resp.Kvs {
			jsonString += fmt.Sprintf(`"%s":"%s"`, ev.Key, ev.Value)
		}
		jsonString = fmt.Sprintf(`{%s}`, jsonString)
		_ = json.Unmarshal([]byte(jsonString), &deploymentNameList)
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
		resp := etcd.Get(etcdClient, key)
		var deploymentNameList []string
		jsonString := ``
		for _, ev := range resp.Kvs {
			jsonString += fmt.Sprintf(`"%s":"%s"`, ev.Key, ev.Value)
		}
		jsonString = fmt.Sprintf(`{%s}`, jsonString)
		_ = json.Unmarshal([]byte(jsonString), &deploymentNameList)
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
	resp := etcd.Get(etcdClient, key)
	var deploymentNameList []string
	jsonString := ``
	for _, ev := range resp.Kvs {
		jsonString += fmt.Sprintf(`"%s":"%s"`, ev.Key, ev.Value)
	}
	jsonString = fmt.Sprintf(`{%s}`, jsonString)
	_ = json.Unmarshal([]byte(jsonString), &deploymentNameList)
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
	resp := etcd.Get(etcdClient, podName)
	pod := def.Pod{}
	jsonString := ``
	for _, ev := range resp.Kvs {
		jsonString += fmt.Sprintf(`"%s":"%s"`, ev.Key, ev.Value)
	}
	jsonString = fmt.Sprintf(`{%s}`, jsonString)
	_ = json.Unmarshal([]byte(jsonString), &pod)
	return &pod
}

func GetPodInstanceResourceUsageByName(etcdClient *clientv3.Client, podInstanceID string) *def.ResourceUsage {
	key := def.GetKeyOfResourceUsageByPodInstanceID(podInstanceID)
	resp := etcd.Get(etcdClient, key)
	resourceUsage := def.ResourceUsage{}
	jsonString := ``
	for _, ev := range resp.Kvs {
		jsonString += fmt.Sprintf(`"%s":"%s"`, ev.Key, ev.Value)
	}
	jsonString = fmt.Sprintf(`{%s}`, jsonString)
	_ = json.Unmarshal([]byte(jsonString), &resourceUsage)
	return &resourceUsage
}
