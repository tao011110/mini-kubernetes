package get_api

import (
	"encoding/json"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"mini-kubernetes/tools/def"
	"mini-kubernetes/tools/etcd"
)

func GetGPUJob(cli *clientv3.Client, jobName string) (def.GPUJobDetail, bool) {
	flag := false
	jobGet := def.GPUJobDetail{
		Job:         def.GPUJob{},
		Pod:         def.Pod{},
		PodInstance: def.PodInstance{},
	}
	replicaIDList := make([]string, 0)
	{
		kv := etcd.Get(cli, def.GetGPUJobKeyByName(jobName)).Kvs
		value := make([]byte, 0)
		if len(kv) != 0 {
			value = kv[0].Value
			flag = true
			err := json.Unmarshal(value, &jobGet.Job)
			if err != nil {
				fmt.Printf("%v\n", err)
				panic(err)
			}
		}
	}
	if !flag {
		return jobGet, flag
	}
	flag = false
	{
		kv := etcd.Get(cli, def.GetKeyOfPod(jobGet.Job.PodName)).Kvs
		value := make([]byte, 0)
		if len(kv) != 0 {
			value = kv[0].Value
			flag = true
			err := json.Unmarshal(value, &jobGet.Pod)
			if err != nil {
				fmt.Printf("%v\n", err)
				panic(err)
			}
		}
	}
	if !flag {
		return jobGet, true
	}
	{
		kv := etcd.GetWithPrefix(cli, def.GetKeyOfPodReplicasNameListByPodName(jobGet.Job.PodName)).Kvs
		value := make([]byte, 0)
		if len(kv) != 0 {
			value = kv[0].Value
			flag = true
			err := json.Unmarshal(value, &replicaIDList)
			if err != nil {
				fmt.Printf("%v\n", err)
				panic(err)
			}
		}
	}
	if len(replicaIDList) < 1 {
		return jobGet, true
	}
	{
		kv := etcd.Get(cli, replicaIDList[0]).Kvs
		value := make([]byte, 0)
		if len(kv) != 0 {
			value = kv[0].Value
			flag = true
			err := json.Unmarshal(value, &jobGet.PodInstance)
			if err != nil {
				fmt.Printf("%v\n", err)
				panic(err)
			}
		}
	}
	return jobGet, true
}
