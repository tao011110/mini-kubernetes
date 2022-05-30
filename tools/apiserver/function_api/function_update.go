package function_api

import (
	clientv3 "go.etcd.io/etcd/client/v3"
	"mini-kubernetes/tools/apiserver/apiserver_utils"
	"mini-kubernetes/tools/def"
	"mini-kubernetes/tools/functional"
)

// softUpdate不删除已经创建的podInstance, 只会将使用的镜像更新为新的镜像, 通过function的动态伸缩逐渐将podInstance替换为新的
// TODO: 需要确认function确实已经存在旧版本
func SoftUpdateFunction(cli *clientv3.Client, function def.Function) {
	pod, _ := functional.GenerateFunctionPodAndService(&function)
	oldFunction := apiserver_utils.GetFunctionByName(cli, function.Name)
	oldPod := apiserver_utils.GetPodByName(cli, oldFunction.PodName)
	oldPod.Spec.Containers[0].Image = pod.Spec.Containers[0].Image
	apiserver_utils.PersistFunction(cli, function)
	apiserver_utils.PersistPod(cli, oldPod)
}

// HardUpdateFunction softUpdate删除已经创建的podInstance, 即重新创建function
func HardUpdateFunction(cli *clientv3.Client, function def.Function) {
	DeleteFunction(cli, function.Name)
	CreateFunction(cli, function)
}
