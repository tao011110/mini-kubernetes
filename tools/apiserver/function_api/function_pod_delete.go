package function_api

import (
	"encoding/json"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"mini-kubernetes/tools/apiserver/apiserver_utils"
	"mini-kubernetes/tools/apiserver/delete_api"
	"mini-kubernetes/tools/apiserver/get_api"
	"mini-kubernetes/tools/def"
	"mini-kubernetes/tools/etcd"
	"mini-kubernetes/tools/util"
)

// DeleteFuncPodInstance 只删除podInstance而不删除pod, 注意从service中同步删除, 参数是podInstanceID
func DeleteFuncPodInstance(cli *clientv3.Client, podName string) (bool, def.Service) {
	flag := true
	var instanceIDList []string
	instanceIDListkey := def.GetKeyOfPodReplicasNameListByPodName(podName)
	resp := etcd.Get(cli, instanceIDListkey)
	if len(resp.Kvs) == 0 {
		return false, def.Service{}
	}
	util.EtcdUnmarshal(resp, &instanceIDList)
	if len(instanceIDList) == 0 {
		return false, def.Service{}
	}
	podInstanceID := instanceIDList[0]
	podInstance := def.PodInstance{}
	podInstanceValue := etcd.Get(cli, podInstanceID).Kvs[0].Value
	err := json.Unmarshal(podInstanceValue, &podInstance)
	if err != nil {
		fmt.Printf("%v\n", err)
		panic(err)
	}
	//etcd.Delete(cli, podInstanceID)

	//更新PodInstanceIDList
	podInstanceIDList := make([]string, 0)
	tmpList := make([]string, 0)
	kvs := etcd.Get(cli, def.PodInstanceListID).Kvs
	if len(kvs) != 0 {
		podInstanceIDListValue := kvs[0].Value
		err := json.Unmarshal(podInstanceIDListValue, &podInstanceIDList)
		if err != nil {
			fmt.Printf("%v\n", err)
			panic(err)
		}
	}
	for _, podInstanceID := range podInstanceIDList {
		if podInstanceID != podInstance.ID {
			tmpList = append(tmpList, podInstanceID)
		}
	}
	podInstanceIDList = tmpList
	podInstanceIDValue, err := json.Marshal(podInstanceIDList)
	if err != nil {
		fmt.Printf("%v\n", err)
		panic(err)
	}
	etcd.Put(cli, def.PodInstanceListID, string(podInstanceIDValue))

	//更新ReplicasNameList
	newInstanceIDList := instanceIDList[1:]
	instanceIDListValue, _ := json.Marshal(newInstanceIDList)
	etcd.Put(cli, podInstance.ID, string(instanceIDListValue))

	// 在service中删除该podInstance
	service, _ := get_api.GetService(cli, "service_"+podInstance.Pod.Metadata.Name[4:])
	fmt.Println(service)
	service.PortsBindings = delete_api.RemovePodInstanceFromService(podInstance, *service)
	apiserver_utils.PersistService(cli, *service)

	return flag, *service
}
