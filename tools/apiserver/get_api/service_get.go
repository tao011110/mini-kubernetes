package get_api

import (
	"encoding/json"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"mini-kubernetes/tools/def"
	"mini-kubernetes/tools/etcd"
)

func GetService(cli *clientv3.Client, serviceName string) (*def.Service, bool) {
	flag := false
	serviceKey := "/service/" + serviceName
	kv := etcd.Get(cli, serviceKey).Kvs
	service := def.Service{}
	serviceValue := make([]byte, 0)
	if len(kv) != 0 {
		serviceValue = kv[0].Value
		flag = true
		err := json.Unmarshal(serviceValue, &service)
		if err != nil {
			fmt.Printf("%v\n", err)
			panic(err)
		}
	}

	return &service, flag
}
