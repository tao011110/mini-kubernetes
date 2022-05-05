package util

import (
	"encoding/json"
	clientv3 "go.etcd.io/etcd/client/v3"
	"mini-kubernetes/tools/def"
	"mini-kubernetes/tools/etcd"
)

func PersistPodInstance(podInstance def.PodInstance, cli *clientv3.Client) {
	byts, _ := json.Marshal(podInstance)
	etcd.Put(cli, podInstance.ID, string(byts))
}
