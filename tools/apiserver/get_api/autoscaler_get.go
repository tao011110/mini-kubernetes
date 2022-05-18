package get_api

import (
	"encoding/json"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"mini-kubernetes/tools/def"
	"mini-kubernetes/tools/etcd"
)

func GetAutoscaler(cli *clientv3.Client, autoscalerName string) (*def.Autoscaler, bool) {
	flag := false
	autoscalerKey := def.GetKeyOfDeployment(autoscalerName)
	kv := etcd.Get(cli, autoscalerKey).Kvs
	autoscaler := def.Autoscaler{}
	autoscalerValue := make([]byte, 0)
	if len(kv) != 0 {
		autoscalerValue = kv[0].Value
		flag = true
		err := json.Unmarshal(autoscalerValue, &autoscaler)
		if err != nil {
			fmt.Printf("%v\n", err)
			panic(err)
		}
	}

	return &autoscaler, flag
}
