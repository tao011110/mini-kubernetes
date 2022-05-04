package apiserver_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mini-kubernetes/tools/def"
	"mini-kubernetes/tools/httpget"
	"mini-kubernetes/tools/master"
	"mini-kubernetes/tools/pod"
	"mini-kubernetes/tools/yaml"
	"net"
	"testing"
)

func Test(t *testing.T) {
	node := def.Node{}
	node.LocalPort = 80
	node.NodeName = "node1"
	node.MasterIpAndPort = master.IP + ":" + master.Port
	addr, _ := net.InterfaceAddrs()
	for _, address := range addr {
		if ip, flag_ := address.(*net.IPNet); flag_ && !ip.IP.IsLoopback() {
			if ip.IP.To4() != nil {
				node.NodeIP = ip.IP.To4()
			}
		}
	}

	//test register_node
	response := def.RegisterToMasterResponse{}
	request := def.RegisterToMasterRequest{
		NodeName:  node.NodeName,
		LocalIP:   node.NodeIP,
		LocalPort: node.LocalPort,
	}
	body, _ := json.Marshal(request)
	err, status := httpget.Post("http://" + node.MasterIpAndPort + "/register_node").
		ContentType("application/json").
		Body(bytes.NewReader(body)).
		GetJson(&response).
		Execute()
	if err != nil {
		fmt.Println("err")
		fmt.Println(err)
	}
	node.NodeID = response.NodeID
	node.NodeName = response.NodeName
	node.CniIP = response.CniIP
	fmt.Printf("register_node is %s and response is: %v\n", status, response)

	//test create_pod
	//需要发送给apiserver的参数为 pod_ def.Pod
	pod_, _ := yaml.ReadPodYamlConfig("../docker_test/docker_test3.yaml")
	request2 := *pod_
	response2 := ""
	body2, _ := json.Marshal(request2)
	err, status = httpget.Post("http://" + node.MasterIpAndPort + "/create_pod").
		ContentType("application/json").
		Body(bytes.NewReader(body2)).
		GetString(&response2).
		Execute()
	if err != nil {
		fmt.Println("err")
		fmt.Println(err)
	}
	fmt.Printf("create_pod is %s and response is: %s\n", status, response2)

	//test get_pod
	//需要发送给apiserver的参数为 podName string
	podName := "pod3"
	response3 := pod.Pod{}
	err, status = httpget.Get("http://" + node.MasterIpAndPort + "/get_pod/" + podName).
		ContentType("application/json").
		GetJson(&response3).
		Execute()
	if err != nil {
		fmt.Println("err")
		fmt.Println(err)
	}
	fmt.Printf("get_pod status is %s\n", status)
	if status == "200" {
		fmt.Printf("get pod_ %s successfully and the response is: %v\n", podName, response3)
	} else {
		fmt.Printf("pod_ %s doesn't exist\n", podName)
	}

	//test delete_pod
	//需要发送给apiserver的参数为 podName string
	response4 := ""
	err, status = httpget.DELETE("http://" + node.MasterIpAndPort + "/delete_pod/" + podName).
		ContentType("application/json").
		GetString(&response4).
		Execute()
	if err != nil {
		fmt.Println("err")
		fmt.Println(err)
	}

	fmt.Printf("get_pod status is %s\n", status)
	if status == "200" {
		fmt.Printf("delete pod_ %s successfully and the response is: %v\n", podName, response4)
	} else {
		fmt.Printf("pod_ %s doesn't exist\n", podName)
	}
}
