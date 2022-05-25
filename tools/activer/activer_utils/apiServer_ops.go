package activer_utils

import clientv3 "go.etcd.io/etcd/client/v3"

//TODO: apiserver加接口, service元数据预存在etcd中但不部署, 只需通过name部署和删除(此处删除不删除元数据)

//func StartService(serviceName string) {
//	//TO DO: apiServer start the service
//}
//
//func StopService(serviceName string) {
//	//TO DO: apiserver stop the service, **but keep the meta**
//}

func AddNPodInstance(podName string, num int) {
	//TODO: apiServer add a podInstance
}

func RemovePodInstance(podName string, num int) {
	//TODO: apiServer delete a podInstance
}

func AdjustReplicaNum2Target(etcdClient *clientv3.Client, funcName string, target int) {
	function := GetFunctionByName(etcdClient, funcName)
	replicaNameList := GetPodReplicaIDListByPodName(etcdClient, function.PodName)
	if len(replicaNameList) < target {
		AddNPodInstance(function.PodName, target-len(replicaNameList))
	} else if len(replicaNameList) > target {
		RemovePodInstance(function.PodName, len(replicaNameList)-target)
		//if target == 0 {
		//	StopService(function.ServiceName)
		//}
	}
}
