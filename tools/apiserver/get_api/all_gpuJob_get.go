package get_api

import (
	"encoding/json"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"mini-kubernetes/tools/def"
	"mini-kubernetes/tools/etcd"
)

func GetAllGPUJob(cli *clientv3.Client) ([]def.GPUJobDetail, bool) {
	dnsPrefix := "/gpu_job/"
	kvs := etcd.GetWithPrefix(cli, dnsPrefix).Kvs
	value := make([]byte, 0)
	jobDetailList := make([]def.GPUJobDetail, 0)
	flag := false
	if len(kvs) != 0 {
		flag = true
		for _, kv := range kvs {
			job := def.GPUJob{}
			value = kv.Value
			err := json.Unmarshal(value, &job)
			if err != nil {
				fmt.Printf("%v\n", err)
				panic(err)
			}
			gpuJobDetail, _ := GetGPUJob(cli, job.Name)
			jobDetailList = append(jobDetailList, gpuJobDetail)
		}
	}
	return jobDetailList, flag
}
