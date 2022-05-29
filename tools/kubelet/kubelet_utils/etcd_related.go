package kubelet_utils

import (
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
