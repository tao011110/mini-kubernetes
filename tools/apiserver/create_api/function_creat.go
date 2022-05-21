package create_api

import (
	clientv3 "go.etcd.io/etcd/client/v3"
	"mini-kubernetes/tools/def"
	"mini-kubernetes/tools/functional"
)

func CreatFunction(cli *clientv3.Client, function def.Function) {
	functional.MakeFunctionalImage(&function)
	pod_, service := functional.GenerateFunctionPodAndService(function)
	function.PodName = pod_.Metadata.Name
	function.ServiceName = service.Metadata.Name
	//持久化pod信息
	{

	}
	//持久化
}
