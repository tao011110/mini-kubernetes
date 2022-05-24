package get_api

import (
	"encoding/json"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"mini-kubernetes/tools/def"
	"mini-kubernetes/tools/etcd"
)

func GetAllStateMachine(cli *clientv3.Client) ([]def.StateMachine, bool) {
	prefix := "/state_machine/"
	kvs := etcd.GetWithPrefix(cli, prefix).Kvs
	value := make([]byte, 0)
	list := make([]def.StateMachine, 0)
	flag := false
	if len(kvs) != 0 {
		flag = true
		for _, kv := range kvs {
			stateMachine := def.StateMachine{}
			value = kv.Value
			err := json.Unmarshal(value, &stateMachine)
			if err != nil {
				fmt.Printf("%v\n", err)
				panic(err)
			}
			list = append(list, stateMachine)
		}
	}

	return list, flag
}
