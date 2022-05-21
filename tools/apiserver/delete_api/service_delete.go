package delete_api

import (
	"encoding/json"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"mini-kubernetes/tools/def"
	"mini-kubernetes/tools/etcd"
)

func DeleteService(cli *clientv3.Client, serviceName string) (string, bool, string) {
	// 在etcd中删除service
	serviceKey := "/service/" + serviceName
	clusterIP := ""
	resp := etcd.Get(cli, serviceKey)
	fmt.Println(serviceKey + "  is serviceKey")
	if len(resp.Kvs) == 0 {
		return clusterIP, false, ""
	}
	etcd.Delete(cli, serviceKey)

	serviceValue := resp.Kvs[0].Value
	service := def.Service{}
	err := json.Unmarshal(serviceValue, &service)
	clusterIP = service.ClusterIP
	if err != nil {
		fmt.Printf("%v\n", err)
		panic(err)
	}

	return clusterIP, true, service.Type
}
