package apiserver_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mini-kubernetes/tools/def"
	"mini-kubernetes/tools/httpget"
	"mini-kubernetes/tools/master"
	"mini-kubernetes/tools/yaml"
	"net"
	"testing"
)

func Test(t *testing.T) {
	node := def.Node{}
	node.LocalPort = 80
	node.NodeName = "node1"
	node.MasterIpAndPort = master.Ip + ":" + master.Port
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
	err := httpget.Post("http://" + node.MasterIpAndPort + "/register_node").
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
	fmt.Printf("register_node response is: %v\n", response)

	//test create_pod
	pod, _ := yaml.ReadYamlConfig("../docker_test/docker_test3.yaml")
	request2 := *pod
	response2 := ""
	body2, _ := json.Marshal(request2)
	err = httpget.Post("http://" + node.MasterIpAndPort + "/create_pod").
		ContentType("application/json").
		Body(bytes.NewReader(body2)).
		GetString(&response2).
		Execute()
	if err != nil {
		fmt.Println("err")
		fmt.Println(err)
	}
	fmt.Printf("create_pod response is: %s\n", response2)

	//test get_pod
	podName := "pod3"
	response3 := def.Pod{}
	err = httpget.Get("http://" + node.MasterIpAndPort + "/get_pod/" + podName).
		ContentType("application/json").
		GetJson(&response3).
		Execute()
	if err != nil {
		fmt.Println("err")
		fmt.Println(err)
	}

	fmt.Printf("get_pod response is: %v\n", response3)

	//test delete_pod
	response4 := ""
	err = httpget.DELETE("http://" + node.MasterIpAndPort + "/delete_pod/" + podName).
		ContentType("application/json").
		GetString(&response4).
		Execute()
	if err != nil {
		fmt.Println("err")
		fmt.Println(err)
	}
	fmt.Printf("delete_pod response is: %v\n", response4)
}
