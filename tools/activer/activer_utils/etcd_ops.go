package activer_utils

import (
	clientv3 "go.etcd.io/etcd/client/v3"
	"mini-kubernetes/tools/def"
	"mini-kubernetes/tools/etcd"
	"mini-kubernetes/tools/util"
)

func GetFunctionNameList(etcdClient *clientv3.Client) []string {
	var functionNameList []string
	util.EtcdUnmarshal(etcd.Get(etcdClient, def.FunctionNameListKey), &functionNameList)
	return functionNameList
}

func GetFunctionByName(etcdClient *clientv3.Client, name string) *def.Function {
	function := def.Function{}
	util.EtcdUnmarshal(etcd.Get(etcdClient, def.GetKeyOfFunction(name)), &function)
	return &function
}

func GetStateMachineByName(etcdClient *clientv3.Client, name string) *def.StateMachine {
	stateMachine := def.StateMachine{}
	util.EtcdUnmarshal(etcd.Get(etcdClient, def.GetKeyOfStateMachine(name)), &stateMachine)
	return &stateMachine
}

func GetPodReplicaIDListByPodName(etcdClient *clientv3.Client, podName string) []string {
	key := def.GetKeyOfPodReplicasNameListByPodName(podName)
	var instanceIDList []string
	util.EtcdUnmarshal(etcd.Get(etcdClient, key), &instanceIDList)
	return instanceIDList
}

func GetServiceByName(etcdClient *clientv3.Client, serviceName string) *def.Service {
	key := def.GetKeyOfService(serviceName)
	service := def.Service{}
	util.EtcdUnmarshal(etcd.Get(etcdClient, key), &service)
	return &service
}
