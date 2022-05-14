package util

import (
	"encoding/json"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"mini-kubernetes/tools/def"
	"mini-kubernetes/tools/etcd"
)

func GetPodInstance(podInstanceName string, cli *clientv3.Client) def.PodInstance {
	resp := etcd.Get(cli, podInstanceName)
	podInstance := def.PodInstance{}
	jsonString := ``
	for _, ev := range resp.Kvs {
		jsonString += fmt.Sprintf(`"%s":"%s", `, ev.Key, ev.Value)
	}
	jsonString = fmt.Sprintf(`{%s}`, jsonString)
	_ = json.Unmarshal([]byte(jsonString), &podInstance)
	return podInstance
}
