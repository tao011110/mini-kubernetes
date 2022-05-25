package function_api

import (
	clientv3 "go.etcd.io/etcd/client/v3"
	"mini-kubernetes/tools/apiserver/apiserver_utils"
	"mini-kubernetes/tools/apiserver/create_api"
	"mini-kubernetes/tools/apiserver/get_api"
	"mini-kubernetes/tools/def"
	"mini-kubernetes/tools/util"
)

func CreateFuncPodInstance(cli *clientv3.Client, podName string, num int) ([]def.PodInstance, def.Service) {
	pod_ := apiserver_utils.GetPodByPodName(cli, podName)
	podInstance := def.PodInstance{}
	podInstance.Pod = pod_

	//将新创建的podInstance写入到etcd当中
	podInstanceKey := def.GetKeyOfPodInstance(pod_.Metadata.Name)
	podInstance.ID = podInstanceKey
	podInstance.ContainerSpec = make([]def.ContainerStatus, len(pod_.Spec.Containers))

	util.PersistPodInstance(podInstance, cli)
	apiserver_utils.AddPodInstanceIDToList(cli, podInstance.ID)

	// 在service中加上该podInstance
	service, _ := get_api.GetService(cli, "service_"+podName[:4])
	service.PortsBindings = create_api.AddPodInstanceIntoService(podInstance, *service)
	apiserver_utils.PersistService(cli, *service)

	return podInstance, *service
}
