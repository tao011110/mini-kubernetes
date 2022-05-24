package apiserver_utils

import (
	"encoding/json"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"mini-kubernetes/tools/def"
	"mini-kubernetes/tools/etcd"
	"mini-kubernetes/tools/util"
)

func GetGPUJobByName(li *clientv3.Client, jobName string) def.GPUJob {
	gpuJob := def.GPUJob{}
	key := def.GetGPUJobKeyByName(jobName)
	util.EtcdUnmarshal(etcd.Get(li, key), &gpuJob)
	return gpuJob
}

func GetPodReplicaListByPodName(li *clientv3.Client, podName string) []string {
	var list []string
	key := def.GetKeyOfPodReplicasNameListByPodName(podName)
	util.EtcdUnmarshal(etcd.Get(li, key), &list)
	return list
}

func GetPodInstanceByID(li *clientv3.Client, id string) def.PodInstance {
	instance := def.PodInstance{}
	util.EtcdUnmarshal(etcd.Get(li, id), &instance)
	return instance
}

func PersistPod(li *clientv3.Client, pod_ def.Pod) {
	key := def.GetKeyOfPod(pod_.Metadata.Name)
	value, err := json.Marshal(pod_)
	if err != nil {
		fmt.Printf("%v\n", err)
		panic(err)
	}
	etcd.Put(li, key, string(value))
}

func GetPodByPodName(li *clientv3.Client, podName string) def.Pod {
	pod_ := def.Pod{}
	key := def.GetKeyOfPod(podName)
	util.EtcdUnmarshal(etcd.Get(li, key), &pod_)
	return pod_
}

func GetPodReplicaIDListByPodName(li *clientv3.Client, podName string) []string {
	idList := make([]string, 0)
	key := def.GetKeyOfPodReplicasNameListByPodName(podName)
	util.EtcdUnmarshal(etcd.Get(li, key), &idList)
	return idList
}

func PersistStateMachine(li *clientv3.Client, stateMachine def.StateMachine) {
	key := def.GetKeyOfStateMachine(stateMachine.Name)
	value, err := json.Marshal(stateMachine)
	if err != nil {
		fmt.Printf("%v\n", err)
		panic(err)
	}
	etcd.Put(li, key, string(value))
}

func PersistService(li *clientv3.Client, service def.Service) {
	key := def.GetKeyOfService(service.Name)
	value, err := json.Marshal(service)
	if err != nil {
		fmt.Printf("%v\n", err)
		panic(err)
	}
	etcd.Put(li, key, string(value))
}

func PersistGPUJob(li *clientv3.Client, job def.GPUJob) {
	key := def.GetGPUJobKeyByName(job.Name)
	value, err := json.Marshal(job)
	if err != nil {
		fmt.Printf("%v\n", err)
		panic(err)
	}
	etcd.Put(li, key, string(value))
}

func PersistFunction(li *clientv3.Client, function def.Function) {
	key := def.GetKeyOfFunction(function.Name)
	value, err := json.Marshal(function)
	if err != nil {
		fmt.Printf("%v\n", err)
		panic(err)
	}
	etcd.Put(li, key, string(value))
}

func AddFunctionNameToList(li *clientv3.Client, functionName string) {
	key := def.FunctionNameListKey
	var list []string
	util.EtcdUnmarshal(etcd.Get(li, key), &list)
	list = append(list, functionName)
	value, err := json.Marshal(list)
	if err != nil {
		fmt.Printf("%v\n", err)
		panic(err)
	}
	etcd.Put(li, key, string(value))
}

func AddPodInstanceIDToList(li *clientv3.Client, id string) {
	key := def.PodInstanceListID
	var list []string
	util.EtcdUnmarshal(etcd.Get(li, key), &list)
	list = append(list, id)
	value, err := json.Marshal(list)
	if err != nil {
		fmt.Printf("%v\n", err)
		panic(err)
	}
	etcd.Put(li, key, string(value))
}
