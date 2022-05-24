package function_api

import (
	clientv3 "go.etcd.io/etcd/client/v3"
	"mini-kubernetes/tools/apiserver/apiserver_utils"
	"mini-kubernetes/tools/def"
	"mini-kubernetes/tools/util"
)

//TODO:

func CreatePodInstance(cli *clientv3.Client, podName string) def.PodInstance {
	pod_ := apiserver_utils.GetPodByPodName(cli, podName)
	podInstance := def.PodInstance{}
	podInstance.Pod = pod_

	//将新创建的podInstance写入到etcd当中
	podInstanceKey := def.GetKeyOfPodInstance(pod_.Metadata.Name)
	podInstance.ID = podInstanceKey
	podInstance.ContainerSpec = make([]def.ContainerStatus, len(pod_.Spec.Containers))

	util.PersistPodInstance(podInstance, cli)
	apiserver_utils.AddPodInstanceIDToList(cli, podInstance.ID)
	return podInstance
}
