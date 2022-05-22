package apiserver_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mini-kubernetes/tools/def"
	"mini-kubernetes/tools/httpget"
	"mini-kubernetes/tools/util"
	"mini-kubernetes/tools/yaml"
	"net"
	"testing"
	"time"
)

var node = def.Node{
	LocalPort:       80,
	ProxyPort:       3000,
	NodeName:        "node1",
	MasterIpAndPort: util.GetLocalIP().String() + ":" + fmt.Sprintf("%d", def.MasterPort),
}

func testRegisterNode() {
	//测试时需修改为本机IP
	node.NodeIP = net.IPv4(192, 168, 1, 7)

	//test register_node
	response := def.RegisterToMasterResponse{}
	request := def.RegisterToMasterRequest{
		NodeName:  node.NodeName,
		LocalIP:   node.NodeIP,
		LocalPort: node.LocalPort,
		ProxyPort: node.ProxyPort,
	}
	fmt.Println("node.MasterIpAndPort is " + node.MasterIpAndPort)
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
	//testGetAllPod()

	testGetAllPodStatus()
}

//TODO: 用来创建clusterIP service，需要发送给apiserver的参数为 service_c  (def.ClusterIPSvc)
func testCreateCIService(path string) {
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

//TODO: 用来创建nodeport service，需要发送给apiserver的参数为 serviceN (def.NodePortSvc)
func testCreateNPService(path string) {
	serviceN, _ := yaml.ReadServiceNodeportConfig(path)
	request2 := *serviceN
	response2 := ""
	body2, _ := json.Marshal(request2)
	err, status := httpget.Post("http://" + node.MasterIpAndPort + "/create/nodePortService").
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

//TODO: 用来删除service，需要发送给apiserver的参数为 serviceName(string)
func testDeleteService() {
	//需要发送给apiserver的参数为 serviceName string
	serviceName := "test-service"
	response4 := ""
	err, status := httpget.DELETE("http://" + node.MasterIpAndPort + "/delete/service/" + serviceName).
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

//TODO: 用来获取特定名称的 service，需要发送给apiserver的参数为 serviceName(string)
func testGetService(serviceName string) {
	//http调用返回的json需解析转为def.Service类型，
	//该结构体的字段，满足了与  K8S中的kubectl describe service serviceName操作返回内容中  所有的信息
	response := def.Service{}
	err, status := httpget.Get("http://" + node.MasterIpAndPort + "/get/service/" + serviceName).
		ContentType("application/json").
		GetJson(&response).
		Execute()
	if err != nil {
		fmt.Println("err")
		fmt.Println(err)
	}
	fmt.Printf("get_service status is %s\n", status)
	if status == "200" {
		fmt.Printf("get service %s successfully and the response is: %v\n", serviceName, response)
	} else {
		fmt.Printf("service %s doesn't exist\n", serviceName)
	}
}

//TODO: 用来获取所有的 service
func testGetAllService() {
	//http调用返回的json需解析转为[]def.Service类型，
	//def.Service 结构体的字段，满足了与  K8S中的kubectl get service操作返回内容中  所有的信息
	//但需要注意的是，我这里提供的是StartTime(time.Time), kubectl在获取之后需要使用如下操作计算出AGE：
	//t := time.Now()  用于获取当前时间
	//Age := t.Sub(podInstance.StartTime)  进行计算，得到AGE
	response := make([]def.Service, 0)
	err, status := httpget.Get("http://" + node.MasterIpAndPort + "/get/all/service").
		ContentType("application/json").
		GetJson(&response).
		Execute()
	if err != nil {
		fmt.Println("err")
		fmt.Println(err)
	}
	fmt.Printf("get_all_service status is %s\n", status)
	if status == "200" {
		fmt.Println("All services' information is as follows")
		for _, service := range response {
			fmt.Printf("%v\n", service)
		}
	} else {
		fmt.Printf("No service exists\n")
	}
}

//TODO: 用来创建DNS和Gateway
func testCreateDNSAndGateway(path string) {
	//需要发送给apiserver的参数为 dns def.DNS
	dns, _ := yaml.ReadDNSConfig(path)
	request2 := *dns
	response2 := ""
	body2, _ := json.Marshal(request2)
	err, status := httpget.Post("http://" + node.MasterIpAndPort + "/create/dns").
		ContentType("application/json").
		Body(bytes.NewReader(body2)).
		GetString(&response2).
		Execute()
	if err != nil {
		fmt.Println("err")
		fmt.Println(err)
	}
	fmt.Printf("create_gateway is %s and response is: %s\n", status, response2)
}

//TODO: 用来获取特定名称的 dns，需要发送给apiserver的参数为 dnsName(string)
func testGetDNS(dnsName string) {
	//http调用返回的json需解析转为def.DNSDetail类型，
	response := def.DNSDetail{}
	err, status := httpget.Get("http://" + node.MasterIpAndPort + "/get/dns/" + dnsName).
		ContentType("application/json").
		GetJson(&response).
		Execute()
	if err != nil {
		fmt.Println("err")
		fmt.Println(err)
	}
	fmt.Printf("get_dns status is %s\n", status)
	if status == "200" {
		fmt.Printf("get dns %s successfully and the response is: %v\n", dnsName, response)
	} else {
		fmt.Printf("dns %s doesn't exist\n", dnsName)
	}
}

//TODO: 用来获取所有的 dns
//我会将所有的DNSDetail都返回，kubectl反馈给用户可以只取每个DNSDetail当中给的一部分信息,
//就像kubectl get service只展示每个service的部分信息
func testGetAllDNS() {
	//http调用返回的json需解析转为[]def.DNSDetail类型
	response := make([]def.DNSDetail, 0)
	err, status := httpget.Get("http://" + node.MasterIpAndPort + "/get/all/dns").
		ContentType("application/json").
		GetJson(&response).
		Execute()
	if err != nil {
		fmt.Println("err")
		fmt.Println(err)
	}
	fmt.Printf("get_all_dns status is %s\n", status)
	if status == "200" {
		fmt.Println("All dns' information is as follows")
		for _, service := range response {
			fmt.Printf("%v\n", service)
		}
	} else {
		fmt.Printf("No dns exists\n")
	}
}

func TestUpdateIptablesRule(t *testing.T) {
	//testRegisterNode()
	//time.Sleep(10 * time.Second)
	//path = "./podForService2.yaml"
	//testCreatePod(path)
	//time.Sleep(10 * time.Second)
	//
	path := "./clusterIPService_test.yaml"
	testCreateCIService(path)
	testGetAllService()

	time.Sleep(5 * time.Second)
	path = "./podForService.yaml"
	testCreatePod(path)

	//time.Sleep(5 * time.Second)
	//testDeleteCIService()

	//time.Sleep(15 * time.Second)
	//path = "./nodePortService_test.yaml"
	//testCreateNPService(path)

	time.Sleep(15 * time.Second)
	//testGetService("test-service2")
	//
	testGetAllService()

	//time.Sleep(5 * time.Second)
	//testDeleteNPService()
}

func TestDNSAndGateway(t *testing.T) {
	path := "./dns_test.yaml"
	testCreateDNSAndGateway(path)

	testGetDNS("dns-name")
	testGetAllDNS()
}
