package apiserver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	clientv3 "go.etcd.io/etcd/client/v3"
	"mini-kubernetes/tools/apiserver/create_api"
	"mini-kubernetes/tools/apiserver/delete_api"
	"mini-kubernetes/tools/apiserver/get_api"
	"mini-kubernetes/tools/apiserver/register_api"
	"mini-kubernetes/tools/def"
	"mini-kubernetes/tools/httpget"
	"strconv"
)

var IpAndPort string
var cli *clientv3.Client

func Start(masterIp string, port string, client *clientv3.Client) {
	IpAndPort = masterIp + ":" + port
	cli = client

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// handle register-api
	e.POST("/register_node", handleRegisterNode)

	// handle create-api
	e.POST("/create_pod", handleCreatePod)
	e.POST("/create/clusterIPService", handleCreateClusterIPService)
	e.POST("/create/nodePortService", handleCreateNodePortService)

	// handle delete-api
	e.DELETE("/delete_pod/:podName", handleDeletePod)
	e.DELETE("/delete/clusterIPService/:serviceName", handleDeleteClusterIPService)
	e.DELETE("/delete/nodePortService/:serviceName", handleDeleteNodePortService)

	// handle get-api
	e.GET("/get_pod/:podName", handleGetPod)
	e.GET("/get/all/pod", handleGetAllPod)
	e.GET("/get/all/podStatus", handleGetAllPodStatus)
	e.GET("/get/all/service", handleGetAllService)
	e.GET("/get/service/:serviceName", handleGetService)

	e.Logger.Fatal(e.Start(":" + port))
}

func handleRegisterNode(c echo.Context) error {
	request := def.RegisterToMasterRequest{}
	response := def.RegisterToMasterResponse{}

	requestBody := new(bytes.Buffer)
	_, err := requestBody.ReadFrom(c.Request().Body)
	if err != nil {
		fmt.Printf("%v\n", err)
		panic(err)
	}

	err = json.Unmarshal(requestBody.Bytes(), &request)
	if err != nil {
		fmt.Printf("%v\n", err)
		panic(err)
	}

	//进行注册
	nodeID, cniIP := register_api.RegisterNode(cli, request, IpAndPort)

	//返回节点编号和节点名称
	response.NodeID = nodeID
	response.NodeName = request.NodeName
	response.CniIP = cniIP
	fmt.Println("Node has registered")

	return c.JSON(200, response)
}

func handleCreatePod(c echo.Context) error {
	pod_ := def.Pod{}
	requestBody := new(bytes.Buffer)
	_, err := requestBody.ReadFrom(c.Request().Body)
	if err != nil {
		fmt.Printf("%v\n", err)
		panic(err)
	}

	err = json.Unmarshal(requestBody.Bytes(), &pod_)
	if err != nil {
		fmt.Printf("%v\n", err)
		panic(err)
	}

	nodeID := create_api.CreatePod(cli, pod_)
	fmt.Println("Pod " + pod_.Metadata.Name + " has been created at node " + strconv.Itoa(nodeID))

	return c.String(200, "Pod "+pod_.Metadata.Name+" has been created at node "+strconv.Itoa(nodeID))
}

func handleCreateClusterIPService(c echo.Context) error {
	service_c := def.ClusterIPSvc{}
	requestBody := new(bytes.Buffer)
	_, err := requestBody.ReadFrom(c.Request().Body)
	if err != nil {
		fmt.Printf("%v\n", err)
		panic(err)
	}

	err = json.Unmarshal(requestBody.Bytes(), &service_c)
	if err != nil {
		fmt.Printf("%v\n", err)
		panic(err)
	}

	service := create_api.CreateClusterIPService(cli, service_c)
	fmt.Println("Service " + service.Name)

	// 创建携程告知所有node上的kube-proxy，使得正在处理的http请求可以立即返回
	nodeList := get_api.GetAllNode(cli)
	for _, node := range nodeList {
		go letProxyCreateCIRule(service, node)
	}

	return c.String(200, "Service "+service.Name)
}

func handleCreateNodePortService(c echo.Context) error {
	service_n := def.NodePortSvc{}
	requestBody := new(bytes.Buffer)
	_, err := requestBody.ReadFrom(c.Request().Body)
	if err != nil {
		fmt.Printf("%v\n", err)
		panic(err)
	}

	err = json.Unmarshal(requestBody.Bytes(), &service_n)
	if err != nil {
		fmt.Printf("%v\n", err)
		panic(err)
	}

	service := create_api.CreateNodePortService(cli, service_n)
	fmt.Println("Service " + service.Name)

	// 创建携程告知所有node上的kube-proxy，使得正在处理的http请求可以立即返回
	nodeList := get_api.GetAllNode(cli)
	for _, node := range nodeList {
		go letProxyCreateNPRule(service, node)
	}

	return c.String(200, "Service "+service.Name)
}

func letProxyCreateCIRule(service def.Service, node def.Node) {
	// 更新所有node的kube-proxy
	target := node.NodeIP.String() + ":" + strconv.Itoa(node.ProxyPort)

	// 创建携程，并发执行
	go func(target string) {
		fmt.Println("target is " + target)
		response := ""
		body, _ := json.Marshal(service)
		err, _ := httpget.Post("http://" + target + "/add/clusterIPServiceRule").
			ContentType("application/json").
			Body(bytes.NewReader(body)).
			GetString(&response).
			Execute()
		if err != nil {
			fmt.Println("err")
			fmt.Println(err)
		}
		fmt.Printf("%s create service successfully\n", target)
	}(target)
}

func letProxyCreateNPRule(service def.Service, node def.Node) {
	// 更新所有node的kube-proxy
	target := node.NodeIP.String() + ":" + strconv.Itoa(node.ProxyPort)

	// 创建携程，并发执行
	go func(target string) {
		fmt.Println("target is " + target)
		response := ""
		body, _ := json.Marshal(service)
		err, _ := httpget.Post("http://" + target + "/add/nodePortServiceRule").
			ContentType("application/json").
			Body(bytes.NewReader(body)).
			GetString(&response).
			Execute()
		if err != nil {
			fmt.Println("err")
			fmt.Println(err)
		}
		fmt.Printf("%s create service successfully\n", target)
	}(target)
}

func handleDeletePod(c echo.Context) error {
	podName := c.Param("podName")

	if delete_api.DeletePod(cli, podName) == true {
		fmt.Println("Pod " + podName + " has been deleted")
		return c.String(200, "Pod "+podName+" has been deleted")
	} else {
		fmt.Println("Pod " + podName + " has been deleted")
		return c.String(404, "Pod "+podName+" doesn't exist")
	}
}

func handleDeleteClusterIPService(c echo.Context) error {
	serviceName := c.Param("serviceName")

	clusterIP, flag := delete_api.DeleteService(cli, serviceName)
	if flag == true {
		// 创建携程告知所有node上的kube-proxy，使得正在处理的http请求可以立即返回
		nodeList := get_api.GetAllNode(cli)
		for _, node := range nodeList {
			go letProxyDeleteCIRule(clusterIP, node)
		}
		fmt.Println("Service " + serviceName + " has been deleted")
		return c.String(200, "Service "+serviceName+" has been deleted")
	} else {
		fmt.Println("Service " + serviceName + " has been deleted")
		return c.String(404, "Service "+serviceName+" doesn't exist")
	}
}

func handleDeleteNodePortService(c echo.Context) error {
	serviceName := c.Param("serviceName")

	clusterIP, flag := delete_api.DeleteService(cli, serviceName)
	if flag == true {
		// 创建携程告知所有node上的kube-proxy，使得正在处理的http请求可以立即返回
		nodeList := get_api.GetAllNode(cli)
		for _, node := range nodeList {
			go letProxyDeleteNPRule(clusterIP, node)
		}
		fmt.Println("Service " + serviceName + " has been deleted")
		return c.String(200, "Service "+serviceName+" has been deleted")
	} else {
		fmt.Println("Service " + serviceName + " has been deleted")
		return c.String(404, "Service "+serviceName+" doesn't exist")
	}
}

func letProxyDeleteCIRule(clusterIP string, node def.Node) {
	// 更新所有node的kube-proxy
	target := node.NodeIP.String() + ":" + strconv.Itoa(node.ProxyPort)
	fmt.Println("target is " + target)

	// 创建携程，并发执行
	go func(target string) {
		response := ""
		err, status := httpget.DELETE("http://" + target + "/delete/clusterIPServiceRule/" + clusterIP).
			ContentType("application/json").
			GetString(&response).
			Execute()
		if err != nil {
			fmt.Println("err")
			fmt.Println(err)
		}

		fmt.Printf("get_pod status is %s\n", status)
		if status == "200" {
			fmt.Printf("%s delete service rule %s successfully\n", target, clusterIP)
		} else {
			fmt.Printf("%s failed to delete service %s\n", target, clusterIP)
		}
	}(target)
}

func letProxyDeleteNPRule(clusterIP string, node def.Node) {
	// 更新所有node的kube-proxy
	target := node.NodeIP.String() + ":" + strconv.Itoa(node.ProxyPort)
	fmt.Println("target is " + target)

	// 创建携程，并发执行
	go func(target string) {
		response := ""
		err, status := httpget.DELETE("http://" + target + "/delete/nodePortServiceRule/" + clusterIP).
			ContentType("application/json").
			GetString(&response).
			Execute()
		if err != nil {
			fmt.Println("err")
			fmt.Println(err)
		}

		fmt.Printf("get_pod status is %s\n", status)
		if status == "200" {
			fmt.Printf("%s delete service rule %s successfully\n", target, clusterIP)
		} else {
			fmt.Printf("%s failed to delete service %s\n", target, clusterIP)
		}
	}(target)
}

func handleGetPod(c echo.Context) error {
	podName := c.Param("podName")
	podInstance, flag := get_api.GetPod(cli, podName)
	fmt.Println(podInstance)

	if flag == false {
		return c.JSON(404, podInstance)
	}

	return c.JSON(200, podInstance)
}

func handleGetAllPod(c echo.Context) error {
	fmt.Println("handleGetAllPod")
	podInstanceList, flag := get_api.GetAllPodInstance(cli)
	podInstanceNameList := make([]string, 0)
	for _, podInstance := range podInstanceList {
		podInstanceNameList = append(podInstanceNameList, podInstance.Metadata.Name)
	}
	fmt.Println(podInstanceNameList)

	if flag == false {
		return c.JSON(404, podInstanceNameList)
	}

	return c.JSON(200, podInstanceNameList)
}

func handleGetAllPodStatus(c echo.Context) error {
	fmt.Println("handleGetAllPodStatus")
	podInstanceBriefList, flag := get_api.GetAllPodInstanceStatus(cli)
	fmt.Println(podInstanceBriefList)

	if flag == false {
		return c.JSON(404, podInstanceBriefList)
	}

	return c.JSON(200, podInstanceBriefList)
}

func handleGetService(c echo.Context) error {
	serviceName := c.Param("serviceName")
	service, flag := get_api.GetService(cli, serviceName)
	fmt.Println(service)

	if flag == false {
		return c.JSON(404, service)
	}

	return c.JSON(200, service)
}

func handleGetAllService(c echo.Context) error {
	fmt.Println("handleGetAllPod")
	serviceList, flag := get_api.GetAllService(cli)

	if flag == false {
		return c.JSON(404, serviceList)
	}

	return c.JSON(200, serviceList)
}
