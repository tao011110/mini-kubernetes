package function_api

import (
	"encoding/json"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"mini-kubernetes/tools/apiserver/apiserver_utils"
	"mini-kubernetes/tools/apiserver/create_api"
	"mini-kubernetes/tools/def"
	"mini-kubernetes/tools/etcd"
	"mini-kubernetes/tools/functional"
	"net"
)

func CreateFunction(cli *clientv3.Client, function def.Function) def.Service {
	pod, ciService := functional.GenerateFunctionPodAndService(&function)
	clusterIP := net.IPv4(10, 24, 0, byte(GetRegisteredNodeID(cli)))
	ciService.Spec.ClusterIP = clusterIP.String()

	apiserver_utils.PersistPod(cli, pod) //NOTE: 只存储不创建
	apiserver_utils.PersistFunction(cli, function)
	apiserver_utils.AddFunctionNameToList(cli, function.Name)
	service := create_api.CreateClusterIPService(cli, ciService)

	return service
}

func GetRegisteredNodeID(cli *clientv3.Client) int {
	key := "/Persistence/funcServiceID/"
	kvs := etcd.Get(cli, key).Kvs
	funcServiceID := 1
	if len(kvs) != 0 {
		funcServiceIDValue := kvs[0].Value
		err := json.Unmarshal(funcServiceIDValue, &funcServiceID)
		if err != nil {
			fmt.Printf("%v\n", err)
			panic(err)
		}
	}
	newFuncServiceID := funcServiceID + 1
	newNodeIDValue, err := json.Marshal(newFuncServiceID)
	if err != nil {
		fmt.Printf("%v\n", err)
		panic(err)
	}
	etcd.Put(cli, key, string(newNodeIDValue))

	return funcServiceID
}
