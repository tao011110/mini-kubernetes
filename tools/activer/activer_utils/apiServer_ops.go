package activer_utils

import (
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"mini-kubernetes/tools/util"
)

func AdjustReplicaNum2Target(etcdClient *clientv3.Client, funcName string, target int) {
	function := GetFunctionByName(etcdClient, funcName)
	replicaNameList := GetPodReplicaIDListByPodName(etcdClient, function.PodName)
	fmt.Println("target size is:   ", target)
	fmt.Println("len(replicaNameList) is:   ", len(replicaNameList))
	if len(replicaNameList) < target {
		util.AddNPodInstance(function.PodName, target-len(replicaNameList))
	} else if len(replicaNameList) > target {
		util.RemovePodInstance(function.PodName, len(replicaNameList)-target)
		//if target == 0 {
		//	StopService(function.ServiceName)
		//}
	}
}
