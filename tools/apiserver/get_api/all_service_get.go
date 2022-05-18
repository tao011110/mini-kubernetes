package get_api

import (
	"encoding/json"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"mini-kubernetes/tools/def"
	"mini-kubernetes/tools/etcd"
)

func GetAllService(cli *clientv3.Client) ([]def.Service, bool) {
	flag := false
	servicePrefix := "/service/"
	kvs := etcd.GetWithPrefix(cli, servicePrefix).Kvs
	serviceValue := make([]byte, 0)
	serviceList := make([]def.Service, 0)
	if len(kvs) != 0 {
		for _, kv := range kvs {
			service := def.Service{}
			serviceValue = kv.Value
			err := json.Unmarshal(serviceValue, &service)
			if err != nil {
				fmt.Printf("%v\n", err)
				panic(err)
			}
			fmt.Println("service.Name is " + service.Name)
			serviceList = append(serviceList, service)
		}
		flag = true
	}

	return serviceList, flag
}
