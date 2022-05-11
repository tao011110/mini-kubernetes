package apiserver_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mini-kubernetes/tools/def"
	"mini-kubernetes/tools/docker"
	"mini-kubernetes/tools/httpget"
	"mini-kubernetes/tools/master"
	"mini-kubernetes/tools/pod"
	"mini-kubernetes/tools/yaml"
	"net"
	"testing"
)

var node = def.Node{
	LocalPort:       4000,
	ProxyPort:       3000,
	NodeName:        "node1",
	MasterIpAndPort: master.IP + ":" + master.Port,
}

func testRegisterNode() {
	//测试时需修改为本机IP
	node.NodeIP = net.IPv4(192, 168, 47, 19)

	//test register_node
	response := def.RegisterToMasterResponse{}
	request := def.RegisterToMasterRequest{
		NodeName:  node.NodeName,
		LocalIP:   node.NodeIP,
		LocalPort: node.LocalPort,
		ProxyPort: node.ProxyPort,
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

	docker.CreateNetBridge("10.24.1.0")
}

func testCreatePod(path string) {
	//需要发送给apiserver的参数为 pod_ def.Pod
	pod_, _ := yaml.ReadPodYamlConfig(path)
	request2 := *pod_
	response2 := ""
	body2, _ := json.Marshal(request2)
	err, status := httpget.Post("http://" + node.MasterIpAndPort + "/create_pod").
		ContentType("application/json").
		Body(bytes.NewReader(body2)).
		GetString(&response2).
		Execute()
	if err != nil {
		fmt.Println("err")
		fmt.Println(err)
	}
	fmt.Printf("create_pod is %s and response is: %s\n", status, response2)

	//TODO:在kubelet正常运行后，这部分测试代码可以删除
	podInstance := def.PodInstance{}
	podInstance.Pod = *pod_
	podInstance.NodeID = uint64(node.NodeID)
	cniIP := net.IPv4(10, 24, 0, 0)
	node.CniIP = net.IP(cniIP)
	podInstance.ID = "/podInstance/" + pod_.Metadata.Name
	podInstance.ContainerSpec = make([]def.ContainerStatus, len(pod_.Spec.Containers))
	pod.CreateAndStartPod(&podInstance, &node)
}

func testGetPod() {
	//需要发送给apiserver的参数为 podName string
	podName := "pod3"
	response3 := def.Pod{}
	err, status := httpget.Get("http://" + node.MasterIpAndPort + "/get_pod/" + podName).
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
}

func testDeletePod() {
	//需要发送给apiserver的参数为 podName string
	podName := "pod3"
	response4 := ""
	err, status := httpget.DELETE("http://" + node.MasterIpAndPort + "/delete_pod/" + podName).
		ContentType("application/json").
		GetString(&response4).
		Execute()
	if err != nil {
		fmt.Println("err")
		fmt.Println(err)
	}

	fmt.Printf("delete_pod status is %s\n", status)
	if status == "200" {
		fmt.Printf("delete pod_ %s successfully and the response is: %v\n", podName, response4)
	} else {
		fmt.Printf("pod_ %s doesn't exist\n", podName)
	}
}

func testGetAllPod() {
	response5 := make([]string, 0)
	err, status := httpget.Get("http://" + node.MasterIpAndPort + "/get/all/pod").
		ContentType("application/json").
		GetJson(&response5).
		Execute()
	if err != nil {
		fmt.Println("err")
		fmt.Println(err)
	}
	fmt.Printf("get_all_pod status is %s\n", status)
	if status == "200" {
		fmt.Println("All pods are as follows")
		for _, podInstance := range response5 {
			fmt.Printf("%v\n", podInstance)
		}
	} else {
		fmt.Printf("No pod exists\n")
	}
}

func testGetAllPodStatus() {
	response5 := make([]def.PodInstanceBrief, 0)
	err, status := httpget.Get("http://" + node.MasterIpAndPort + "/get/all/podStatus").
		ContentType("application/json").
		GetJson(&response5).
		Execute()
	if err != nil {
		fmt.Println("err")
		fmt.Println(err)
	}
	fmt.Printf("get_all_pod status is %s\n", status)
	if status == "200" {
		fmt.Println("All pods' brief information is as follows")
		for _, podInstanceBrief := range response5 {
			fmt.Printf("%v\n", podInstanceBrief)
		}
	} else {
		fmt.Printf("No pod exists\n")
	}
}

func TestPod(t *testing.T) {
	testRegisterNode()

	var path = "./podForService.yaml"
	testCreatePod(path)
	//
	//testGetPod()
	//
	//testDeletePod()
	//
	testGetAllPod()

	testGetAllPodStatus()
}

func testCreateCIService(path string) {
	//需要发送给apiserver的参数为 service_c def.ClusterIP
	serviceC, _ := yaml.ReadServiceClusterIPConfig(path)
	request2 := *serviceC
	response2 := ""
	body2, _ := json.Marshal(request2)
	err, status := httpget.Post("http://" + node.MasterIpAndPort + "/create/clusterIPService").
		ContentType("application/json").
		Body(bytes.NewReader(body2)).
		GetString(&response2).
		Execute()
	if err != nil {
		fmt.Println("err")
		fmt.Println(err)
	}
	fmt.Printf("create_service is %s and response is: %s\n", status, response2)
}

func testDeleteCIService() {
	//需要发送给apiserver的参数为 serviceName string
	serviceName := "test-service"
	response4 := ""
	err, status := httpget.DELETE("http://" + node.MasterIpAndPort + "/delete/clusterIPService/" + serviceName).
		ContentType("application/json").
		GetString(&response4).
		Execute()
	if err != nil {
		fmt.Println("err")
		fmt.Println(err)
	}

	fmt.Printf("delete clusterIPService status is %s\n", status)
	if status == "200" {
		fmt.Printf("delete clusterIPService %s successfully and the response is: %v\n", serviceName, response4)
	} else {
		fmt.Printf("clusterIPService %s doesn't exist\n", serviceName)
	}
}

func TestUpdateIptablesRule(t *testing.T) {
	testRegisterNode()

	var path = "./podForService.yaml"
	testCreatePod(path)
	//path = "./podForService2.yaml"
	//testCreatePod(path)

	path = "./clusterIPService_test.yaml"
	testCreateCIService(path)

	testDeleteCIService()
}
