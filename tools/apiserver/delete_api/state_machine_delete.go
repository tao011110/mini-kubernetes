package delete_api

import (
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"mini-kubernetes/tools/def"
	"mini-kubernetes/tools/etcd"
)

func DeleteStateMachine(cli *clientv3.Client, stateMachineName string) bool {
	key := def.GetKeyOfStateMachine(stateMachineName)
	resp := etcd.Get(cli, key)
	fmt.Println(key + "  is key")
	if len(resp.Kvs) == 0 {
		return false
	}
	etcd.Delete(cli, key)

	return true
}
