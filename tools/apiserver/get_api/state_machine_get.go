package get_api

import (
	"encoding/json"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"mini-kubernetes/tools/def"
	"mini-kubernetes/tools/etcd"
)

func GetStateMachine(cli *clientv3.Client, stateMachineName string) (def.StateMachine, bool) {
	flag := false
	key := def.GetKeyOfStateMachine(stateMachineName)
	kv := etcd.Get(cli, key).Kvs
	stateMachine := def.StateMachine{}
	dnsValue := make([]byte, 0)
	if len(kv) != 0 {
		dnsValue = kv[0].Value
		flag = true
		err := json.Unmarshal(dnsValue, &stateMachine)
		if err != nil {
			fmt.Printf("%v\n", err)
			panic(err)
		}
	}

	return stateMachine, flag
}
