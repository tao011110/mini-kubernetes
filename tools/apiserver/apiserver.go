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

	// handle delete-api
	e.DELETE("/delete_pod/:podName", handleDeletePod)

	// handle get-api
	e.GET("/get_pod/:podName", handleGetPod)

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
	pod := def.Pod{}
	requestBody := new(bytes.Buffer)
	_, err := requestBody.ReadFrom(c.Request().Body)
	if err != nil {
		fmt.Printf("%v\n", err)
		panic(err)
	}

	err = json.Unmarshal(requestBody.Bytes(), &pod)
	if err != nil {
		fmt.Printf("%v\n", err)
		panic(err)
	}

	nodeID := create_api.CreatePod(cli, pod)
	fmt.Println("Pod " + pod.Metadata.Name + " has been created at node " + strconv.Itoa(nodeID))

	return c.String(200, "Pod "+pod.Metadata.Name+" has been created at node "+strconv.Itoa(nodeID))
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

func handleGetPod(c echo.Context) error {
	podName := c.Param("podName")
	podInstance, flag := get_api.GetPod(cli, podName)
	fmt.Println(podInstance)

	if flag == false {
		return c.JSON(404, podInstance)
	}

	return c.JSON(200, podInstance)
}
