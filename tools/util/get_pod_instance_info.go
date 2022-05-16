package util

import (
	"encoding/json"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"mini-kubernetes/tools/def"
	"mini-kubernetes/tools/etcd"
)

func GetPodInstance(podInstanceName string, cli *clientv3.Client) def.PodInstance {
	kv := etcd.Get(cli, podInstanceName).Kvs
	podInstance := def.PodInstance{}
	podInstanceValue := make([]byte, 0)
	if len(kv) != 0 {
		podInstanceValue = kv[0].Value
		err := json.Unmarshal(podInstanceValue, &podInstance)
		if err != nil {
			fmt.Printf("%v\n", err)
			panic(err)
		}
	}

	return podInstance
}
