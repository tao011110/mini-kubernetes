package get_api

import (
	"encoding/json"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"mini-kubernetes/tools/def"
	"mini-kubernetes/tools/etcd"
)

func GetAllPodInstance(cli *clientv3.Client) ([]def.PodInstance, bool) {
	flag := false
	podInstanceKey := "/podInstance/"
	kvs := etcd.GetWithPrefix(cli, podInstanceKey).Kvs
	podInstanceValue := make([]byte, 0)
	podInstanceList := make([]def.PodInstance, 0)
	if len(kvs) != 0 {
		for _, kv := range kvs {
			podInstance := def.PodInstance{}
			podInstanceValue = kv.Value
			err := json.Unmarshal(podInstanceValue, &podInstance)
			if err != nil {
				fmt.Printf("%v\n", err)
				panic(err)
			}
			fmt.Println("podInstance.Metadata.Name is " + podInstance.Metadata.Name)
			podInstanceList = append(podInstanceList, podInstance)
		}
		flag = true
	}

	return podInstanceList, flag
}
