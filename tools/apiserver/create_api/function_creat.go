package create_api

import (
	clientv3 "go.etcd.io/etcd/client/v3"
	"mini-kubernetes/tools/apiserver/apiserver_utils"
	"mini-kubernetes/tools/def"
	"mini-kubernetes/tools/functional"
)

func CreateFunction(cli *clientv3.Client, function def.Function) def.Service {
	pod, ciService := functional.GenerateFunctionPodAndService(&function)
	apiserver_utils.PersistPod(cli, pod) //NOTE: 只存储不创建
	apiserver_utils.PersistFunction(cli, function)
	apiserver_utils.AddFunctionNameToList(cli, function.Name)
	service := CreateClusterIPService(cli, ciService)

	return service
}
