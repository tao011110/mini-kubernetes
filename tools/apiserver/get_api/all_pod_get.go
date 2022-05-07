package get_api

import (
	"encoding/json"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"mini-kubernetes/tools/def"
	"mini-kubernetes/tools/etcd"
)

func GetAllPod(cli *clientv3.Client) ([]string, bool) {
	flag := false
	podInstanceKey := "/podInstance/"
	kvs := etcd.GetWithPrefix(cli, podInstanceKey).Kvs
	podInstance := def.PodInstance{}
	podInstanceValue := make([]byte, 0)
	podInstanceNameList := make([]string, 0)
	if len(kvs) != 0 {
		for _, kv := range kvs {
			podInstanceValue = kv.Value
			err := json.Unmarshal(podInstanceValue, &podInstance)
			if err != nil {
				fmt.Printf("%v\n", err)
				panic(err)
			}
			fmt.Println("podInstance.Metadata.Name is " + podInstance.Metadata.Name)
			podInstanceNameList = append(podInstanceNameList, podInstance.Metadata.Name)
		}
		flag = true
	}

	return podInstanceNameList, flag
}
