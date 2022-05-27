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

// 用来创建Deployment，需要发送给apiserver的参数为 deployment (def.Deployment)
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

// 用来删除deployment，需要发送给apiserver的参数为 deploymentName(string)
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

// 用来获取特定名称的 deployment，需要发送给apiserver的参数为 deploymentName(string)
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

// 用来获取所有的 deployment
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

// 用来创建AutoScaler，需要发送给apiserver的参数为 autoScaler (def.AutoScaler)
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

// 用来删除autoscaler，需要发送给apiserver的参数为 autoscalerName(string)
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

// 用来获取特定名称的 autoscaler，需要发送给apiserver的参数为 autoscalerName(string)
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

// 用来获取所有的 autoscaler
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

func TestReplicas(t *testing.T) {
	//path := "./deployment_test.yaml"
	//testCreateDeployment(path)
	//
	//time.Sleep(30 * time.Second)
	//testGetDeployment("test-deployment")
	//
	//testGetAllDeployment()
	//time.Sleep(10 * time.Second)

	//testDeleteDeployment()

	path := "./HorizontalPodAutoscaler_test.yaml"
	testCreateAutoscaler(path)

	time.Sleep(30 * time.Second)
	testGetAutoscaler("test-autoscaler")

	testGetAllAutoscaler()

	//time.Sleep(30 * time.Second)
	//testDeleteAutoscaler()
}

//TODO: 用来创建GPUJob，需要发送给apiserver的参数为 gpu (def.GPUJob)
func testCreateGPUJob(path string) {
	gpu, _ := yaml.ReadGPUJobConfig(path)
	request2 := *gpu
	response2 := ""
	body2, _ := json.Marshal(request2)
	err, status := httpget.Post("http://" + node.MasterIpAndPort + "/create/gpuJob").
		ContentType("application/json").
		Body(bytes.NewReader(body2)).
		GetString(&response2).
		Execute()
	if err != nil {
		fmt.Println("err")
		fmt.Println(err)
	}
	fmt.Printf("create_gpuJob is %s and response is: %s\n", status, response2)
}

//TODO: 用来获取特定名称的 gpuJob，需要发送给apiserver的参数为 gpuJobName(string)
func testGetGPUJob(gpuJobName string) {
	response := def.GPUJobDetail{}
	err, status := httpget.Get("http://" + node.MasterIpAndPort + "/get/gpuJob/" + gpuJobName).
		ContentType("application/json").
		GetJson(&response).
		Execute()
	if err != nil {
		fmt.Println("err")
		fmt.Println(err)
	}
	fmt.Printf("get_gpuJob status is %s\n", status)
	if status == "200" {
		fmt.Printf("get gpuJob %s successfully and the response is: %v\n", gpuJobName, response)
	} else {
		fmt.Printf("gpuJob %s doesn't exist\n", gpuJobName)
	}
}

//TODO: 用来获取所有的 gpuJob
func testGetAllGPUJob() {
	response := make([]def.GPUJobDetail, 0)
	err, status := httpget.Get("http://" + node.MasterIpAndPort + "/get/all/gpuJob").
		ContentType("application/json").
		GetJson(&response).
		Execute()
	if err != nil {
		fmt.Println("err")
		fmt.Println(err)
	}
	fmt.Printf("get_all_gpuJob status is %s\n", status)
	if status == "200" {
		fmt.Println("All gpuJobs' information is as follows")
		for _, gpuJobDetail := range response {
			fmt.Printf("%v\n", gpuJobDetail)
		}
	} else {
		fmt.Printf("No gpuJob exists\n")
	}
}

func TestGPU(t *testing.T) {
	testCreateGPUJob("./gpu_test.yaml")

	time.Sleep(50 * time.Second)

	testGetGPUJob("GPUJob-test")
	testGetAllGPUJob()
}

//TODO: 用来创建function，需要发送给apiserver的参数为 function (def.Function)
//需要注意的是，返回的 response 会提供一个url，kubectl将它呈现给用户，后续用户可以使用它发送请求
func testCreateFunction(path string) {
	function, _ := yaml.ReadFunctionConfig(path)
	request := *function
	response := ""
	body2, _ := json.Marshal(request)
	err, status := httpget.Post("http://" + node.MasterIpAndPort + "/create/function").
		ContentType("application/json").
		Body(bytes.NewReader(body2)).
		GetString(&response).
		Execute()
	if err != nil {
		fmt.Println("err")
		fmt.Println(err)
	}
	fmt.Printf("create_function is %s and response is: %s\n", status, response)
}

//TODO: 用来获取特定名称的 function，需要发送给apiserver的参数为 functionName(string)
func testGetFunction(functionName string) {
	response := def.Function{}
	err, status := httpget.Get("http://" + node.MasterIpAndPort + "/get/function/" + functionName).
		ContentType("application/json").
		GetJson(&response).
		Execute()
	if err != nil {
		fmt.Println("err")
		fmt.Println(err)
	}
	fmt.Printf("get_function status is %s\n", status)
	if status == "200" {
		fmt.Printf("get function %s successfully and the response is: %v\n", functionName, response)
	} else {
		fmt.Printf("function %s doesn't exist\n", functionName)
	}
}

//TODO: 用来获取所有的 function
func testGetAllFunction() {
	response := make([]def.Function, 0)
	err, status := httpget.Get("http://" + node.MasterIpAndPort + "/get/all/function").
		ContentType("application/json").
		GetJson(&response).
		Execute()
	if err != nil {
		fmt.Println("err")
		fmt.Println(err)
	}
	fmt.Printf("get_all_function status is %s\n", status)
	if status == "200" {
		fmt.Println("All functions' information is as follows")
		for _, function := range response {
			fmt.Printf("%v\n", function)
		}
	} else {
		fmt.Printf("No function exists\n")
	}
}

//该函数不需要加入到kubectl里面
func testCreateFuncPodInstance(podName string) {
	request2 := podName
	response2 := ""
	body2, _ := json.Marshal(request2)
	err, status := httpget.Post("http://" + node.MasterIpAndPort + "/create/funcPodInstance").
		ContentType("application/json").
		Body(bytes.NewReader(body2)).
		GetString(&response2).
		Execute()
	if err != nil {
		fmt.Println("err")
		fmt.Println(err)
	}
	fmt.Printf("create_function is %s and response is: %s\n", status, response2)
}

//该函数不需要加入到kubectl里面
func testDeleteFuncPodInstance(podName string) {
	response4 := ""
	err, status := httpget.DELETE("http://" + node.MasterIpAndPort + "/delete/funcPodInstance/" + podName).
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

func TestFunction(t *testing.T) {
	testCreateFunction("./function_test.yaml")

	testGetFunction("function_test")
	testGetAllFunction()

	//testCreateFuncPodInstance("pod_functional_function_test_1_57cb853e-a6fb-460f-a24a-c2a9b0753d8a")
	//testDeleteFuncPodInstance("pod_functional_function_test_1_57cb853e-a6fb-460f-a24a-c2a9b0753d8a")
}

func TestActiver(t *testing.T) {
	response := ""
	functionName := "function_test"
	type person struct {
		UserType int `json:"userType"`
	}
	request2 := person{
		UserType: 2,
	}
	body2, _ := json.Marshal(request2)
	fmt.Println(request2)
	fmt.Println(body2)
	fmt.Println(bytes.NewReader(body2))
	err, status := httpget.Get("http://127.0.0.1:3306" + "/function/" + functionName + "?test_param1=0&test_param2=2").
		ContentType("application/json").
		Body(bytes.NewReader(body2)).
		GetString(&response).
		Execute()
	if err != nil {
		fmt.Println("err")
		fmt.Println(err)
	}
	fmt.Printf("get_function status is %s\n", status)
	if status == "200" {
		fmt.Printf("get function %s successfully and the response is: %v\n", functionName, response)
	} else {
		fmt.Printf("function %s doesn't exist\n", functionName)
	}
}

//TODO: 用来创建StateMachine，需要发送给apiserver的参数为 stateMachine (def.StateMachine)
//需要注意的是，返回的 response 会提供一个url，kubectl将它呈现给用户，后续用户可以使用它发送请求
func testCreateStateMachine(path string) {
	//TODO: 这里需要注意的是，用户会交给kubectl一个.json文件的路径，kubectl从中json字符串后，将其转为def.StateMachine。
	//这里并没有关于读文件的操作，kubectl需自己实现。所用的json文件，可以参考/presentation/serverless-stateMachine/state_machine.json
	stateMachineJson := "{\n  \"Name\": \"test_state_machine\",\n  \"StartAt\": \"State1\",\n  \"States\": {\n    \"State1\": {\n      \"Type\": \"Task\",\n      \"Resource\": \"state1_function\",\n      \"Next\": \"Choice\"\n    },\n    \"Choice\": {\n      \"Type\": \"Choice\",\n      \"Choices\": [\n        {\n          \"Variable\": \"$.type_state1\",\n          \"StringEquals\": \"1\",\n          \"Next\": \"State2\"\n        },\n        {\n          \"Variable\": \"$.type_state1\",\n          \"StringEquals\": \"2\",\n          \"Next\": \"State3\"\n        }\n      ]\n    },\n    \"State2\": {\n      \"Type\": \"Task\",\n      \"Resource\": \"state2_function\",\n      \"Next\": \"State4\"\n    },\n    \"State3\": {\n      \"Type\": \"Task\",\n      \"Resource\": \"state3_function\",\n      \"Next\": \"State4\"\n    },\n    \"State4\": {\n      \"Type\": \"Task\",\n      \"Resource\": \"state4_function\",\n      \"End\": true\n    }\n  }\n}\n"
	stateMachine := def.StateMachine{}
	_ = json.Unmarshal([]byte(stateMachineJson), &stateMachine)
	fmt.Println(stateMachine)

	request2 := stateMachine
	response2 := ""
	body2, _ := json.Marshal(request2)
	err, status := httpget.Post("http://" + node.MasterIpAndPort + "/create/stateMachine").
		ContentType("application/json").
		Body(bytes.NewReader(body2)).
		GetString(&response2).
		Execute()
	if err != nil {
		fmt.Println("err")
		fmt.Println(err)
	}
	fmt.Printf("create_stateMachine is %s and response is: %s\n", status, response2)
}

//TODO: 用来获取特定名称的 StateMachine，需要发送给apiserver的参数为 stateMachineName(string)
func testGetStateMachine(stateMachineName string) {
	response := def.StateMachine{}
	err, status := httpget.Get("http://" + node.MasterIpAndPort + "/get/stateMachine/" + stateMachineName).
		ContentType("application/json").
		GetJson(&response).
		Execute()
	if err != nil {
		fmt.Println("err")
		fmt.Println(err)
	}
	fmt.Printf("get_stateMachine status is %s\n", status)
	if status == "200" {
		fmt.Printf("get stateMachine %s successfully and the response is: %v\n", stateMachineName, response)
	} else {
		fmt.Printf("stateMachine %s doesn't exist\n", stateMachineName)
	}
}

//TODO: 用来获取所有的 StateMachine
func testGetAllStateMachine() {
	response := make([]def.StateMachine, 0)
	err, status := httpget.Get("http://" + node.MasterIpAndPort + "/get/all/stateMachine").
		ContentType("application/json").
		GetJson(&response).
		Execute()
	if err != nil {
		fmt.Println("err")
		fmt.Println(err)
	}
	fmt.Printf("get_all_stateMachine status is %s\n", status)
	if status == "200" {
		fmt.Println("All stateMachines' information is as follows")
		for _, stateMachine := range response {
			fmt.Printf("%v\n", stateMachine)
		}
	} else {
		fmt.Printf("No stateMachine exists\n")
	}
}

func TestStateStateMachine(t *testing.T) {
	testCreateFunction("./state1.yaml")
	testCreateFunction("./state2.yaml")
	testCreateFunction("./state3.yaml")
	testCreateFunction("./state4.yaml")
	testCreateStateMachine("./function_test.json")

	testGetStateMachine("test_state_machine")
	testGetAllStateMachine()
}

func TestStateActiver(t *testing.T) {
	response := ""
	stateMachineName := "test_state_machine"
	type Person struct {
		Type float64 `json:"type"`
	}
	request2 := Person{
		Type: 1,
	}
	body2, _ := json.Marshal(request2)
	fmt.Println(request2)
	fmt.Println(body2)
	fmt.Println(bytes.NewReader(body2))
	err, status := httpget.Get("http://127.0.0.1:3306" + "/state_machine/" + stateMachineName + "?test_param1=0&test_param2=2").
		ContentType("application/json").
		Body(bytes.NewReader(body2)).
		GetString(&response).
		Execute()
	if err != nil {
		fmt.Println("err")
		fmt.Println(err)
	}
	fmt.Printf("get_function status is %s\n", status)
	if status == "200" {
		fmt.Printf("get stateMachine %s successfully and the response is: %v\n", stateMachineName, response)
	} else {
		fmt.Printf("stateMachine %s doesn't exist\n", stateMachineName)
	}
}
