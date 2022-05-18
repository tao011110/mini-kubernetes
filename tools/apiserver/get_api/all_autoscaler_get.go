package get_api

import (
	"encoding/json"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"mini-kubernetes/tools/def"
	"mini-kubernetes/tools/etcd"
)

func GetAllAutoscaler(cli *clientv3.Client) ([]def.Autoscaler, bool) {
	autoscalerKey := "/autoscaler/"
	kvs := etcd.GetWithPrefix(cli, autoscalerKey).Kvs
	autoscalerValue := make([]byte, 0)
	autoscalerList := make([]def.Autoscaler, 0)
	if len(kvs) != 0 {
		for _, kv := range kvs {
			autoscaler := def.Autoscaler{}
			autoscalerValue = kv.Value
			err := json.Unmarshal(autoscalerValue, &autoscaler)
			if err != nil {
				fmt.Printf("%v\n", err)
				panic(err)
			}
			fmt.Println("autoscaler.Metadata.Name is " + autoscaler.Metadata.Name)
			autoscalerList = append(autoscalerList, autoscaler)
		}
	} else {
		return autoscalerList, false
	}

	return autoscalerList, true
}
