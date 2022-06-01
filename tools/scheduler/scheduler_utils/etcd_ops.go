package scheduler_utils

import (
	"encoding/json"
	clientv3 "go.etcd.io/etcd/client/v3"
	"mini-kubernetes/tools/def"
	"mini-kubernetes/tools/etcd"
	"mini-kubernetes/tools/util"
)

func GetAllPodInstancesOfANode(nodeID int, etcdClient *clientv3.Client) []string {
	var replicaNameList []string
	util.EtcdUnmarshal(etcd.Get(etcdClient, def.PodInstanceListKeyOfNodeID(nodeID)), &replicaNameList)
	return replicaNameList
}

func GetAllPodInstancesID(etcdClient *clientv3.Client) []string {
	var allReplicas []string
	util.EtcdUnmarshal(etcd.Get(etcdClient, def.PodInstanceListID), &allReplicas)
	return allReplicas
}

func GetPodInstanceByName(etcdClient *clientv3.Client, replicaName string) def.PodInstance {
	podInstance := def.PodInstance{}
	util.EtcdUnmarshal(etcd.Get(etcdClient, replicaName), &podInstance)
	return podInstance
}

func GetResourceUsageSequenceByNodeID(etcdClient *clientv3.Client, nodeID int) def.ResourceUsage {
	// TODO: 注册node时添加空ResourceUsage(valid = false)
	nodeResource := def.ResourceUsage{}
	util.EtcdUnmarshal(etcd.Get(etcdClient, def.KeyNodeResourceUsage(nodeID)), &nodeResource)
	return nodeResource
}

func GetAllNodesID(etcdClient *clientv3.Client) []int {
	var nodeIDList []int
	util.EtcdUnmarshal(etcd.Get(etcdClient, def.NodeListID), &nodeIDList)
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
	instance.NodeID = nodeID
	util.PersistPodInstance(*instance, etcdClient)
	replicaNameList := GetAllPodInstancesOfANode(nodeID, etcdClient)
	replicaNameList = append(replicaNameList, instance.ID)
	PersistPodInstanceListOfNode(etcdClient, replicaNameList, nodeID)
}

func PersistPodInstanceListOfNode(etcdClient *clientv3.Client, replicaNameList []string, nodeID int) {
	newJsonString, _ := json.Marshal(replicaNameList)
	etcd.Put(etcdClient, def.PodInstanceListKeyOfNodeID(nodeID), string(newJsonString))
}

func GetPodInstanceIDListOfNode(etcdClient *clientv3.Client, nodeID int) []string {
	key := def.GetKeyOfPodInstanceListKeyOfNodeByID(nodeID)
	var replicaIDList []string
	util.EtcdUnmarshal(etcd.Get(etcdClient, key), &replicaIDList)
	return replicaIDList
}
