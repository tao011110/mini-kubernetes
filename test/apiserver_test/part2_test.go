package apiserver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mini-kubernetes/tools/def"
	"mini-kubernetes/tools/httpget"
	"mini-kubernetes/tools/util"
	"mini-kubernetes/tools/yaml"
	"testing"
	"time"
)

var node = def.Node{
	LocalPort:       80,
	ProxyPort:       3000,
	NodeName:        "node1",
	MasterIpAndPort: fmt.Sprintf("%s:%d", util.GetLocalIP().String(), def.MasterPort),
}

//TODO: 用来创建Deployment，需要发送给apiserver的参数为 deployment (def.Deployment)
func testCreateDeployment(path string) {
	deployment, _ := yaml.ReadDeploymentConfig(path)
	request2 := *deployment
	response2 := ""
	body2, _ := json.Marshal(request2)
	err, status := httpget.Post("http://" + node.MasterIpAndPort + "/create/deployment").
		ContentType("application/json").
		Body(bytes.NewReader(body2)).
		GetString(&response2).
		Execute()
	if err != nil {
		fmt.Println("err")
		fmt.Println(err)
	}
	fmt.Printf("create_deployment is %s and response is: %s\n", status, response2)
}

//TODO: 用来删除deployment，需要发送给apiserver的参数为 deploymentName(string)
func testDeleteDeployment() {
	deploymentName := "test-deployment"
	response4 := ""
	err, status := httpget.DELETE("http://" + node.MasterIpAndPort + "/delete/deployment/" + deploymentName).
		ContentType("application/json").
		GetString(&response4).
		Execute()
	if err != nil {
		fmt.Println("err")
		fmt.Println(err)
	}

	fmt.Printf("delete deployment status is %s\n", status)
	if status == "200" {
		fmt.Printf("delete deployment %s successfully and the response is: %v\n", deploymentName, response4)
	} else {
		fmt.Printf("deployment %s doesn't exist\n", deploymentName)
	}
}

//TODO: 用来获取特定名称的 deployment，需要发送给apiserver的参数为 deploymentName(string)
func testGetDeployment(deploymentName string) {
	//http调用返回的json需解析转为def.DeploymentDetail类型，
	response := def.DeploymentDetail{}
	err, status := httpget.Get("http://" + node.MasterIpAndPort + "/get/deployment/" + deploymentName).
		ContentType("application/json").
		GetJson(&response).
		Execute()
	if err != nil {
		fmt.Println("err")
		fmt.Println(err)
	}
	fmt.Printf("get_deployment status is %s\n", status)
	if status == "200" {
		fmt.Printf("get deployment %s successfully and the response is: %v\n", deploymentName, response)
	} else {
		fmt.Printf("deployment %s doesn't exist\n", deploymentName)
	}
}

//TODO: 用来获取所有的 deployment
func testGetAllDeployment() {
	// DeploymentBrief提供了 的 kubelet get deployment 显示的全部信息
	response := make([]def.DeploymentBrief, 0)
	err, status := httpget.Get("http://" + node.MasterIpAndPort + "/get/all/deployment").
		ContentType("application/json").
		GetJson(&response).
		Execute()
	if err != nil {
		fmt.Println("err")
		fmt.Println(err)
	}
	fmt.Printf("get_all_deployment status is %s\n", status)
	if status == "200" {
		fmt.Println("All deployments' information is as follows")
		for _, deployment := range response {
			fmt.Printf("%v\n", deployment)
		}
	} else {
		fmt.Printf("No deployment exists\n")
	}
}

//TODO: 用来创建AutoScaler，需要发送给apiserver的参数为 autoScaler (def.AutoScaler)
func testCreateAutoscaler(path string) {
	autoscaler, _ := yaml.ReadAutoScalerConfig(path)
	request2 := *autoscaler
	response2 := ""
	body2, _ := json.Marshal(request2)
	err, status := httpget.Post("http://" + node.MasterIpAndPort + "/create/autoscaler").
		ContentType("application/json").
		Body(bytes.NewReader(body2)).
		GetString(&response2).
		Execute()
	if err != nil {
		fmt.Println("err")
		fmt.Println(err)
	}
	fmt.Printf("create_deployment is %s and response is: %s\n", status, response2)
}

//TODO: 用来删除autoscaler，需要发送给apiserver的参数为 autoscalerName(string)
func testDeleteAutoscaler() {
	autoscalerName := "test-autoscaler"
	response4 := ""
	err, status := httpget.DELETE("http://" + node.MasterIpAndPort + "/delete/autoscaler/" + autoscalerName).
		ContentType("application/json").
		GetString(&response4).
		Execute()
	if err != nil {
		fmt.Println("err")
		fmt.Println(err)
	}

	fmt.Printf("delete autoscaler status is %s\n", status)
	if status == "200" {
		fmt.Printf("delete autoscaler %s successfully and the response is: %v\n", autoscalerName, response4)
	} else {
		fmt.Printf("autoscaler %s doesn't exist\n", autoscalerName)
	}
}

//TODO: 用来获取特定名称的 autoscaler，需要发送给apiserver的参数为 autoscalerName(string)
func testGetAutoscaler(autoscalerName string) {
	response := def.AutoscalerDetail{}
	err, status := httpget.Get("http://" + node.MasterIpAndPort + "/get/autoscaler/" + autoscalerName).
		ContentType("application/json").
		GetJson(&response).
		Execute()
	if err != nil {
		fmt.Println("err")
		fmt.Println(err)
	}
	fmt.Printf("get_autoscaler status is %s\n", status)
	if status == "200" {
		fmt.Printf("get autoscaler %s successfully and the response is: %v\n", autoscalerName, response)
	} else {
		fmt.Printf("autoscaler %s doesn't exist\n", autoscalerName)
	}
}

//TODO: 用来获取所有的 autoscaler
func testGetAllAutoscaler() {
	// AutoscalerBrief提供了 的 kubelet get autoscaler 显示的部分信息（根据我们项目与K8S实现的部分差异，对一些信息予以删除）
	response := make([]def.AutoscalerBrief, 0)
	err, status := httpget.Get("http://" + node.MasterIpAndPort + "/get/all/autoscaler").
		ContentType("application/json").
		GetJson(&response).
		Execute()
	if err != nil {
		fmt.Println("err")
		fmt.Println(err)
	}
	fmt.Printf("get_all_autoscaler status is %s\n", status)
	if status == "200" {
		fmt.Println("All autoscalers' information is as follows")
		for _, autoscaler := range response {
			fmt.Printf("%v\n", autoscaler)
		}
	} else {
		fmt.Printf("No autoscaler exists\n")
	}
}

func Test(t *testing.T) {
	//path := "./deployment_test.yaml"
	//testCreateDeployment(path)
	//
	//time.Sleep(30 * time.Second)
	//testGetDeployment("test-deployment")
	//
	//testGetAllDeployment()
	//time.Sleep(10 * time.Second)
	//
	//testDeleteDeployment()

	path := "./HorizontalPodAutoscaler_test.yaml"
	testCreateAutoscaler(path)

	time.Sleep(30 * time.Second)
	testGetAutoscaler("test-autoscaler")

	testGetAllAutoscaler()

	testDeleteAutoscaler()
}
