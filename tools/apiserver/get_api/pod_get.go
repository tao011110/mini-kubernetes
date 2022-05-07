package get_api

import (
	"encoding/json"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"mini-kubernetes/tools/def"
	"mini-kubernetes/tools/etcd"
)

func GetPod(cli *clientv3.Client, podName string) (*def.PodInstance, bool) {
	flag := false
	podInstanceKey := "/podInstance/" + podName
	kv := etcd.Get(cli, podInstanceKey).Kvs
	podInstance := def.PodInstance{}
	podInstanceValue := make([]byte, 0)
	if len(kv) != 0 {
		podInstanceValue = kv[0].Value
		flag = true
		err := json.Unmarshal(podInstanceValue, &podInstance)
		if err != nil {
			fmt.Printf("%v\n", err)
			panic(err)
		}
	}

	return &podInstance, flag
}
