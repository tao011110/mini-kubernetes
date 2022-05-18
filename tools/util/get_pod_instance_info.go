package util

import (
	"encoding/json"
	clientv3 "go.etcd.io/etcd/client/v3"
	"mini-kubernetes/tools/def"
	"mini-kubernetes/tools/etcd"
)

func GetPodInstance(podInstanceID string, cli *clientv3.Client) def.PodInstance {
	resp := etcd.Get(cli, podInstanceID)
	podInstance := def.PodInstance{}
	EtcdUnmarshal(resp, &podInstance)
	return podInstance
}

func EtcdUnmarshal(resp *clientv3.GetResponse, v interface{}) {
	kv := resp.Kvs
	value := make([]byte, 0)
	if len(kv) != 0 {
		value = kv[0].Value
		err := json.Unmarshal(value, v)
		if err != nil {
			panic(err)
		}
	}
}
