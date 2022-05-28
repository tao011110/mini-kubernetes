package function_api

import (
	"encoding/json"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"mini-kubernetes/tools/def"
	"mini-kubernetes/tools/etcd"
)

func GetAllFunction(cli *clientv3.Client) ([]def.Function, bool) {
	prefix := "/function/"
	kvs := etcd.GetWithPrefix(cli, prefix).Kvs
	value := make([]byte, 0)
	list := make([]def.Function, 0)
	flag := false
	if len(kvs) != 0 {
		flag = true
		for _, kv := range kvs {
			function := def.Function{}
			value = kv.Value
			err := json.Unmarshal(value, &function)
			if err != nil {
				fmt.Printf("%v\n", err)
				panic(err)
			}
			list = append(list, function)
		}
	}

	return list, flag
}
