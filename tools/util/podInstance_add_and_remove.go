package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mini-kubernetes/tools/def"
	"mini-kubernetes/tools/httpget"
)

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
		err, status := httpget.Post("http://" + GetLocalIP().String() + ":" + fmt.Sprintf("%d", def.MasterPort) + "/create/replicasPodInstance").
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
		err, status := httpget.DELETE("http://" + GetLocalIP().String() + ":" + fmt.Sprintf("%d", def.MasterPort) + "/delete/replicasPodInstance/" + podName).
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

func RemovePodInstanceByID(podInstanceID string) {
	//apiServer delete a podInstance
	response4 := ""
	err, status := httpget.DELETE("http://" + GetLocalIP().String() + ":" + fmt.Sprintf("%d", def.MasterPort) + "/delete/podInstance/" + podInstanceID).
		ContentType("application/json").
		GetString(&response4).
		Execute()
	if err != nil {
		fmt.Println("err")
		fmt.Println(err)
	}
	fmt.Printf("delete funcPodInstance status is %s\n", status)
	if status == "200" {
		fmt.Printf("delete funcPodInstance %s successfully and the response is: %v\n", podInstanceID, response4)
	} else {
		fmt.Printf("funcPodInstance %s doesn't exist\n", podInstanceID)
	}
}
