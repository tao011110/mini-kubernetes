package scheduler_utils

import (
	"encoding/json"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"mini-kubernetes/tools/def"
	"mini-kubernetes/tools/etcd"
	"mini-kubernetes/tools/util"
)

func GetAllPodInstancesOfANode(nodeID int, etcdClient *clientv3.Client) []string {
	resp := etcd.Get(etcdClient, def.PodInstanceListKeyOfNodeID(nodeID))
	var replicaNameList []string
	jsonString := ``
	for _, ev := range resp.Kvs {
		jsonString += fmt.Sprintf(`"%s":"%s"`, ev.Key, ev.Value)
	}
	jsonString = fmt.Sprintf(`{%s}`, jsonString)
	_ = json.Unmarshal([]byte(jsonString), &replicaNameList)
	return replicaNameList
}

func GetAllPodInstancesID(etcdClient *clientv3.Client) []string {
	resp := etcd.Get(etcdClient, def.PodInstanceListName)
	var allReplicas []string
	jsonString := ``
	for _, ev := range resp.Kvs {
		jsonString += fmt.Sprintf(`"%s":"%s"`, ev.Key, ev.Value)
	}
	jsonString = fmt.Sprintf(`{%s}`, jsonString)
	_ = json.Unmarshal([]byte(jsonString), &allReplicas)
	return allReplicas
}

func GetPodInstanceByName(etcdClient *clientv3.Client, replicaName string) def.PodInstance {
	resp := etcd.Get(etcdClient, replicaName)
	podInstance := def.PodInstance{}
	jsonString := ``
	for _, ev := range resp.Kvs {
		jsonString += fmt.Sprintf(`"%s":"%s", `, ev.Key, ev.Value)
	}
	jsonString = fmt.Sprintf(`{%s}`, jsonString)
	_ = json.Unmarshal([]byte(jsonString), &podInstance)
	return podInstance
}

func GetResourceUsageSequenceByNodeID(etcdClient *clientv3.Client, nodeID int) def.ResourceUsageSequence {
	resp := etcd.Get(etcdClient, def.KeyNodeResourceUsage(nodeID))
	nodeResource := def.ResourceUsageSequence{}
	jsonString := ``
	for _, ev := range resp.Kvs {
		jsonString += fmt.Sprintf(`"%s":"%s", `, ev.Key, ev.Value)
	}
	jsonString = fmt.Sprintf(`{%s}`, jsonString)
	_ = json.Unmarshal([]byte(jsonString), &nodeResource)
	return nodeResource
}

func GetAllNodesID(etcdClient *clientv3.Client) []int {
	resp := etcd.Get(etcdClient, def.NodeListName)
	var nodeIDList []int
	jsonString := ``
	for _, ev := range resp.Kvs {
		jsonString += fmt.Sprintf(`"%s":"%s"`, ev.Key, ev.Value)
	}
	jsonString = fmt.Sprintf(`{%s}`, jsonString)
	_ = json.Unmarshal([]byte(jsonString), &nodeIDList)
	return nodeIDList
}

func DeletePodInstanceFromNode(etcdClient *clientv3.Client, nodeID int, instanceName string) {
	replicaNameList := GetAllPodInstancesOfANode(nodeID, etcdClient)
	for index, replicaName := range replicaNameList {
		if replicaName == instanceName {
			replicaNameList = append(replicaNameList[:index], replicaNameList[index+1:]...)
			break
		}
	}
	PersistPodInstanceListOfNode(etcdClient, replicaNameList, nodeID)
}
func AddPodInstanceToNode(etcdClient *clientv3.Client, nodeID int, instance *def.PodInstance) {
	replicaNameList := GetAllPodInstancesOfANode(nodeID, etcdClient)
	replicaNameList = append(replicaNameList, instance.ID)
	PersistPodInstanceListOfNode(etcdClient, replicaNameList, nodeID)
	instance.NodeID = uint64(nodeID)
	util.PersistPodInstance(*instance, etcdClient)
}

func PersistPodInstanceListOfNode(etcdClient *clientv3.Client, replicaNameList []string, nodeID int) {
	newJsonString, _ := json.Marshal(replicaNameList)
	etcd.Put(etcdClient, def.PodInstanceListKeyOfNodeID(nodeID), string(newJsonString))
}
