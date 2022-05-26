package activer_utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"mini-kubernetes/tools/def"
	"mini-kubernetes/tools/httpget"
	"mini-kubernetes/tools/util"
)

//TODO: apiserver加接口, service元数据预存在etcd中但不部署, 只需通过name部署和删除(此处删除不删除元数据)

//func StartService(serviceName string) {
//	//TO DO: apiServer start the service
//}
//
//func StopService(serviceName string) {
//	//TO DO: apiserver stop the service, **but keep the meta**
//}

func AddNPodInstance(podName string, num int) {
	//apiServer add a podInstance
	for i := 0; i < num; i++ {
		request2 := podName
		response2 := ""
		body2, _ := json.Marshal(request2)
		err, status := httpget.Post("http://" + util.GetLocalIP().String() + ":" + fmt.Sprintf("%d", def.MasterPort) + "/create/funcPodInstance").
			ContentType("application/json").
			Body(bytes.NewReader(body2)).
			GetString(&response2).
			Execute()
		if err != nil {
			fmt.Println("err")
			fmt.Println(err)
		}
		fmt.Printf("create_funcPodInstance is %s and response is: %s\n", status, response2)
	}
}

func RemovePodInstance(podName string, num int) {
	//apiServer delete a podInstance
	for i := 0; i < num; i++ {
		response4 := ""
		err, status := httpget.DELETE("http://" + util.GetLocalIP().String() + ":" + fmt.Sprintf("%d", def.MasterPort) + "/delete/funcPodInstance/" + podName).
			ContentType("application/json").
			GetString(&response4).
			Execute()
		if err != nil {
			fmt.Println("err")
			fmt.Println(err)
		}

		fmt.Printf("delete funcPodInstance status is %s\n", status)
		if status == "200" {
			fmt.Printf("delete funcPodInstance %s successfully and the response is: %v\n", podName, response4)
		} else {
			fmt.Printf("funcPodInstance %s doesn't exist\n", podName)
		}
	}
}

func AdjustReplicaNum2Target(etcdClient *clientv3.Client, funcName string, target int) {
	function := GetFunctionByName(etcdClient, funcName)
	replicaNameList := GetPodReplicaIDListByPodName(etcdClient, function.PodName)
	if len(replicaNameList) < target {
		AddNPodInstance(function.PodName, target-len(replicaNameList))
	} else if len(replicaNameList) > target {
		RemovePodInstance(function.PodName, len(replicaNameList)-target)
		//if target == 0 {
		//	StopService(function.ServiceName)
		//}
	}
}
