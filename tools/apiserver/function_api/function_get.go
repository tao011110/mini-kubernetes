package function_api

import (
	"encoding/json"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"mini-kubernetes/tools/def"
	"mini-kubernetes/tools/etcd"
)

func GetFunction(cli *clientv3.Client, functionName string) (def.Function, bool) {
	flag := false
	key := def.GetKeyOfFunction(functionName)
	kv := etcd.Get(cli, key).Kvs
	function := def.Function{}
	dnsValue := make([]byte, 0)
	if len(kv) != 0 {
		dnsValue = kv[0].Value
		flag = true
		err := json.Unmarshal(dnsValue, &function)
		if err != nil {
			fmt.Printf("%v\n", err)
			panic(err)
		}
	}

	return function, flag
}
