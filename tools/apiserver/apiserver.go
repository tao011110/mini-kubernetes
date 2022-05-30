package apiserver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/robfig/cron"
	clientv3 "go.etcd.io/etcd/client/v3"
	"mini-kubernetes/tools/apiserver/apiserver_utils"
	"mini-kubernetes/tools/apiserver/create_api"
	"mini-kubernetes/tools/apiserver/delete_api"
	"mini-kubernetes/tools/apiserver/function_api"
	"mini-kubernetes/tools/apiserver/get_api"
	"mini-kubernetes/tools/apiserver/gpu_job_api"
	"mini-kubernetes/tools/apiserver/register_api"
	"mini-kubernetes/tools/coredns/coredns_utils"
	"mini-kubernetes/tools/def"
	"mini-kubernetes/tools/etcd"
	"mini-kubernetes/tools/httpget"
	"mini-kubernetes/tools/util"
	"strconv"
	"time"
)

var IpAndPort string
var cli *clientv3.Client
var HeartBeatMap map[int]time.Time
var ShouldStop bool

func Start(masterIp string, port string, client *clientv3.Client) {
	IpAndPort = masterIp + ":" + port
	cli = client
	ShouldStop = false

	//Initialize heartbeat map
	HeartBeatMap = make(map[int]time.Time)
	nodeIDList := apiserver_utils.GetNodeList(client)
	for _, id := range nodeIDList {
		HeartBeatMap[id] = time.Now()
	}
	go NodeChecker()

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// handle register-api
	e.POST("/register_node", handleRegisterNode)

	// handle create-api
	e.POST("/create_pod", handleCreatePod)
	e.POST("/create/clusterIPService", handleCreateClusterIPService)
	e.POST("/create/nodePortService", handleCreateNodePortService)
	e.POST("/create/deployment", handleCreateDeployment)
	e.POST("/create/autoscaler", handleCreateAutoscaler)
	e.POST("/create/dns", handleCreateDNS)
	e.POST("/create/function", handleCreateFunction)
	e.POST("/create/stateMachine", handleCreateStateMachine)
	e.POST("/create/gpuJob", handleCreateGPUJob)
	e.POST("/gpu_job_result", handleOutputGPUJob)
	e.POST("/create/funcPodInstance", handleCreateFuncPodInstance)

	// handle delete-api
	e.DELETE("/delete_pod/:podInstanceName", handleDeletePod)
	e.DELETE("/delete/service/:serviceName", handleDeleteService)
	e.DELETE("/delete/deployment/:deploymentName", handleDeleteDeployment)
	e.DELETE("/delete/autoscaler/:autoscalerName", handleDeleteAutoscaler)
	e.DELETE("/delete/funcPodInstance/:podName", handleDeleteFuncPodInstance)
	e.DELETE("/delete/function/:funcName", handleDeleteFunction)
	e.DELETE("/delete/stateMachine/:stateMachineName", handelDeleteStateMachine)
	e.DELETE("/delete/podInstance/:podInstanceID", handleDeletePodInstance)

	// handle get-api
	e.GET("/get/all/node", handleGetAllNode)
	e.GET("/get_pod/:podInstanceName", handleGetPod)
	e.GET("/get/all/pod", handleGetAllPod)
	e.GET("/get/all/podStatus", handleGetAllPodStatus)
	e.GET("/get/all/service", handleGetAllService)
	e.GET("/get/service/:serviceName", handleGetService)
	e.GET("/get/deployment/:deploymentName", handleGetDeployment)
	e.GET("/get/all/deployment", handleGetAllDeployment)
	e.GET("/get/autoscaler/:autoscalerName", handleGetAutoscaler)
	e.GET("/get/all/autoscaler", handleGetAllAutoscaler)
	e.GET("/get/dns/:dnsName", handleGetDNS)
	e.GET("/get/all/dns", handleGetAllDNS)
	e.GET("/get/gpuJob/:gpuJobName", handleGetGPUJob)
	e.GET("/get/all/gpuJob", handleGetAllGPUJob)
	e.GET("/get/function/:functionName", handleGetFunction)
	e.GET("/get/all/function", handleGetAllFunction)
	e.GET("/get/stateMachine/:stateMachineName", handleGetStateMachine)
	e.GET("/get/all/stateMachine", handleGetAllStateMachine)

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

	HeartBeatMap[nodeID] = time.Now()

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

	_, flag := get_api.GetPodInstance(cli, pod_.Metadata.Name)
	if flag == true {
		return c.String(409, "Pod "+pod_.Metadata.Name+" has already existed")
	}
	podInstance := create_api.CreatePod(cli, pod_)
	fmt.Println("Pod " + pod_.Metadata.Name + " has been created")

	go func(podInstanceID string) {
		fmt.Println("Start to watch ", podInstanceID)
		watchResult := etcd.Watch(cli, podInstanceID)
		for wc := range watchResult {
			change := def.PodInstance{}
			for _, w := range wc.Events {
				if w.Type == clientv3.EventTypePut {
					err := json.Unmarshal(w.Kv.Value, &change)
					if err != nil {
						fmt.Println(err)
						panic(err)
					}
					if change.IP != "" {
						// 创建携程告知所有node上的kube-proxy，使得正在处理的http请求可以立即返回
						serviceList := create_api.CheckAddInService(cli, change)
						nodeList := get_api.GetAllNode(cli)
						for _, service := range serviceList {
							// 需要检验，若原先该IP已存在于service当中，则不再重复添加
							for _, bindings := range service.PortsBindings {
								for _, endPoint := range bindings.Endpoints {
									if endPoint == change.IP {
										return
									}
								}
							}
							if service.Type == "ClusterIP" {
								for _, node := range nodeList {
									letProxyDeleteCIRule(service.ClusterIP, node)
									letProxyCreateCIRule(service, node)
								}
							} else {
								for _, node := range nodeList {
									letProxyDeleteCIRule(service.ClusterIP, node)
									letProxyCreateNPRule(service, node)
								}
							}
						}
						fmt.Println("end watch")
						return
					}
				}
			}
		}
	}(podInstance.ID)

	return c.String(200, "Pod "+pod_.Metadata.Name+" has been created")
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

	_, flag := get_api.GetService(cli, service_c.Metadata.Name)
	if flag == true {
		return c.String(409, "ClusterIPService "+service_c.Metadata.Name+" has already existed")
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

	_, flag := get_api.GetService(cli, service_n.Metadata.Name)
	if flag == true {
		return c.String(409, "NodePortService "+service_n.Metadata.Name+" has already existed")
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

func handleCreateDeployment(c echo.Context) error {
	deployment := def.Deployment{}
	requestBody := new(bytes.Buffer)
	_, err := requestBody.ReadFrom(c.Request().Body)
	if err != nil {
		fmt.Printf("%v\n", err)
		panic(err)
	}
	err = json.Unmarshal(requestBody.Bytes(), &deployment)
	if err != nil {
		fmt.Printf("%v\n", err)
		panic(err)
	}

	_, flag := get_api.GetDeployment(cli, deployment.Metadata.Name)
	if flag == true {
		return c.String(409, "Deployment "+deployment.Metadata.Name+" has already existed")
	}
	create_api.CreateDeployment(cli, deployment)
	fmt.Println("Create deployment ", deployment.Metadata.Name)

	return c.String(200, "deployment "+deployment.Metadata.Name+" has been created")
}

func handleCreateAutoscaler(c echo.Context) error {
	autoscaler := def.Autoscaler{}
	requestBody := new(bytes.Buffer)
	_, err := requestBody.ReadFrom(c.Request().Body)
	if err != nil {
		fmt.Printf("%v\n", err)
		panic(err)
	}
	err = json.Unmarshal(requestBody.Bytes(), &autoscaler)
	if err != nil {
		fmt.Printf("%v\n", err)
		panic(err)
	}

	_, flag := get_api.GetAutoscaler(cli, autoscaler.Metadata.Name)
	if flag == true {
		return c.String(409, "Autoscaler "+autoscaler.Metadata.Name+" has already existed")
	}
	create_api.CreateAutoscaler(cli, autoscaler)
	fmt.Println("Create autoscaler ", autoscaler.Metadata.Name)

	return c.String(200, "autoscaler "+autoscaler.Metadata.Name+" has been created")
}

//TODO: 这里watch可能有问题？
func handleCreateDNS(c echo.Context) error {
	dns := def.DNS{}
	requestBody := new(bytes.Buffer)
	_, err := requestBody.ReadFrom(c.Request().Body)
	if err != nil {
		fmt.Printf("%v\n", err)
		panic(err)
	}
	err = json.Unmarshal(requestBody.Bytes(), &dns)
	if err != nil {
		fmt.Printf("%v\n", err)
		panic(err)
	}

	_, flag := get_api.GetDNS(cli, dns.Name)
	if flag == true {
		return c.String(409, "DNS "+dns.Name+" has already existed")
	}
	dnsDetail, gatewayID := create_api.CreateDNS(cli, dns)
	go func(gatewayID string) {
		fmt.Println("Start to watch ", gatewayID)
		watchResult := etcd.Watch(cli, gatewayID)
		for wc := range watchResult {
			change := def.PodInstance{}
			for _, w := range wc.Events {
				if w.Type == clientv3.EventTypePut {
					err := json.Unmarshal(w.Kv.Value, &change)
					if err != nil {
						fmt.Println(err)
						panic(err)
					}
					if change.IP != "" {
						coredns_utils.AddItem(cli, dnsDetail.Host+":80", change.IP, 80)
						fmt.Println("find add")
					}
				} else {
					if w.Type == clientv3.EventTypeDelete {
						err := json.Unmarshal(w.Kv.Value, &change)
						if err != nil {
							fmt.Println(err)
							panic(err)
						}
						if change.IP != "" {
							coredns_utils.DeleteItem(cli, dnsDetail.Host+":80")
							fmt.Println("find delete")
							return
						}
					}
				}
			}
		}
	}(gatewayID)

	return c.String(200, "DNS "+dns.Name+" has been created")
}

func letProxyCreateCIRule(service def.Service, node def.Node) {
	// 更新所有node的kube-proxy
	target := node.NodeIP.String() + ":" + strconv.Itoa(node.ProxyPort)

	// 创建携程，并发执行
	//go func(target string) {
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
	//}(target)
}

func letProxyCreateNPRule(service def.Service, node def.Node) {
	// 更新所有node的kube-proxy
	target := node.NodeIP.String() + ":" + strconv.Itoa(node.ProxyPort)

	// 创建携程，并发执行
	//go func(target string) {
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
	//}(target)
}

func handleCreateFunction(c echo.Context) error {
	function := def.Function{}
	requestBody := new(bytes.Buffer)
	_, err := requestBody.ReadFrom(c.Request().Body)
	if err != nil {
		fmt.Printf("%v\n", err)
		panic(err)
	}
	err = json.Unmarshal(requestBody.Bytes(), &function)
	if err != nil {
		fmt.Printf("%v\n", err)
		panic(err)
	}

	_, flag := function_api.GetFunction(cli, function.Name)
	if flag == true {
		return c.String(409, "Function "+function.Name+" has already existed")
	}
	function.URL = fmt.Sprintf("http://%s:%d/function/%s", util.GetLocalIP().String(), def.ActiverPort, function.Name)
	service := function_api.CreateFunction(cli, function)
	fmt.Println("Create function ", function.Name)

	// 创建携程告知所有node上的kube-proxy，使得正在处理的http请求可以立即返回
	nodeList := get_api.GetAllNode(cli)
	for _, node := range nodeList {
		go letProxyCreateCIRule(service, node)
	}

	return c.String(200, function.URL)
}

func handleCreateGPUJob(c echo.Context) error {
	job := def.GPUJob{}
	requestBody := new(bytes.Buffer)
	_, err := requestBody.ReadFrom(c.Request().Body)
	if err != nil {
		fmt.Printf("%v\n", err)
		panic(err)
	}
	err = json.Unmarshal(requestBody.Bytes(), &job)
	if err != nil {
		fmt.Printf("%v\n", err)
		panic(err)
	}

	_, flag := get_api.GetGPUJob(cli, job.Name)
	if flag == true {
		return c.String(409, "GPUJob "+job.Name+" has already existed")
	}
	create_api.CreateGPUJobUploader(cli, job)
	fmt.Println("Create job ", job.Name)

	return c.String(200, "job "+job.Name+" has been created")
}

func handleCreateStateMachine(c echo.Context) error {
	stateMachine := def.StateMachine{}
	requestBody := new(bytes.Buffer)
	_, err := requestBody.ReadFrom(c.Request().Body)
	if err != nil {
		fmt.Printf("%v\n", err)
		panic(err)
	}
	err = json.Unmarshal(requestBody.Bytes(), &stateMachine)
	if err != nil {
		fmt.Printf("%v\n", err)
		panic(err)
	}

	_, flag := get_api.GetStateMachine(cli, stateMachine.Name)
	if flag == true {
		return c.String(409, "StateMachine "+stateMachine.Name+" has already existed")
	}
	stateMachine.URL = fmt.Sprintf("http://%s:%d/state_machine/%s", util.GetLocalIP().String(), def.ActiverPort, stateMachine.Name)
	create_api.CreateStateMachine(cli, stateMachine)
	fmt.Println("Create stateMachine ", stateMachine.Name)

	return c.String(200, stateMachine.URL)
}

func handleOutputGPUJob(c echo.Context) error {
	gpuJobResponse := def.GPUJobResponse{}
	requestBody := new(bytes.Buffer)
	_, err := requestBody.ReadFrom(c.Request().Body)
	if err != nil {
		fmt.Printf("%v\n", err)
		panic(err)
	}
	fmt.Println("requestBody.Bytes():  ", string(requestBody.Bytes()))
	err = json.Unmarshal(requestBody.Bytes(), &gpuJobResponse)
	if err != nil {
		fmt.Printf("%v\n", err)
		panic(err)
	}
	gpu_job_api.OutputGPUJOBResponse(cli, gpuJobResponse)
	fmt.Println("Create gpuJobResponse ", gpuJobResponse.JobName)

	return c.String(200, "gpuJobResponse "+gpuJobResponse.JobName+" has been output")
}

func handleCreateFuncPodInstance(c echo.Context) error {
	podName := ""
	requestBody := new(bytes.Buffer)
	_, err := requestBody.ReadFrom(c.Request().Body)
	if err != nil {
		fmt.Printf("%v\n", err)
		panic(err)
	}
	err = json.Unmarshal(requestBody.Bytes(), &podName)
	if err != nil {
		fmt.Printf("%v\n", err)
		panic(err)
	}
	podInstance := function_api.CreateFuncPodInstance(cli, podName)
	fmt.Println("Create FuncPodInstance of ", podInstance.Pod.Metadata.Name)

	go func(podInstanceID string) {
		fmt.Println("Start to watch ", podInstanceID)
		watchResult := etcd.Watch(cli, podInstanceID)
		for wc := range watchResult {
			change := def.PodInstance{}
			for _, w := range wc.Events {
				if w.Type == clientv3.EventTypePut {
					err := json.Unmarshal(w.Kv.Value, &change)
					if err != nil {
						fmt.Println(err)
						panic(err)
					}
					if change.IP != "" {
						fmt.Println("change.IP:  ", change.IP)

						// 在service中加上该podInstance
						fmt.Println("podName:   ", podName)
						service, _ := get_api.GetService(cli, "service_"+podName[4:])
						fmt.Println("and service is", service)
						// 需要检验，若原先该IP已存在于service当中，则不再重复添加
						for _, bindings := range service.PortsBindings {
							for _, endPoint := range bindings.Endpoints {
								if endPoint == change.IP {
									return
								}
							}
						}
						service.PortsBindings = create_api.AddPodInstanceIntoService(change, *service)
						apiserver_utils.PersistService(cli, *service)

						// 创建携程告知所有node上的kube-proxy，使得正在处理的http请求可以立即返回
						nodeList := get_api.GetAllNode(cli)
						for _, node := range nodeList {
							letProxyDeleteCIRule(service.ClusterIP, node)
							letProxyCreateCIRule(*service, node)
						}
						fmt.Println("end watch")
						return
					}
				}
			}
		}
	}(podInstance.ID)

	return c.String(200, "FuncPodInstance of "+podInstance.Pod.Metadata.Name+" has been created")
}

func handleDeletePod(c echo.Context) error {
	podpodInstanceName := c.Param("podInstanceName")
	flag, podInstance := delete_api.DeletePod(cli, podpodInstanceName)
	if flag == true {
		fmt.Println("Pod " + podpodInstanceName + " has been deleted")
		serviceList := delete_api.CheckDeleteInService(cli, podInstance)
		nodeList := get_api.GetAllNode(cli)
		for _, service := range serviceList {
			if service.Type == "ClusterIP" {
				for _, node := range nodeList {
					letProxyDeleteCIRule(service.ClusterIP, node)
					letProxyCreateCIRule(service, node)
				}
			} else {
				for _, node := range nodeList {
					letProxyDeleteCIRule(service.ClusterIP, node)
					letProxyCreateNPRule(service, node)
				}
			}
		}
		return c.String(200, "Pod "+podpodInstanceName+" has been deleted")
	} else {
		fmt.Println("Pod " + podpodInstanceName + " has been deleted")
		return c.String(404, "Pod "+podpodInstanceName+" doesn't exist")
	}
}

func handleDeleteService(c echo.Context) error {
	serviceName := c.Param("serviceName")

	clusterIP, flag, typeName := delete_api.DeleteService(cli, serviceName)
	if flag == true {
		// 创建携程告知所有node上的kube-proxy，使得正在处理的http请求可以立即返回
		nodeList := get_api.GetAllNode(cli)
		if typeName == "ClusterIP" {
			for _, node := range nodeList {
				go letProxyDeleteCIRule(clusterIP, node)
			}
		} else {
			for _, node := range nodeList {
				go letProxyDeleteNPRule(clusterIP, node)
			}
		}
		fmt.Println("Service " + serviceName + " has been deleted")
		return c.String(200, "Service "+serviceName+" has been deleted")
	} else {
		fmt.Println("Service " + serviceName + " has been deleted")
		return c.String(404, "Service "+serviceName+" doesn't exist")
	}
}

func handleDeleteDeployment(c echo.Context) error {
	deploymentName := c.Param("deploymentName")

	if delete_api.DeleteDeployment(cli, deploymentName) == true {
		fmt.Println("Deployment " + deploymentName + " has been deleted")
		return c.String(200, "Deployment "+deploymentName+" has been deleted")
	} else {
		fmt.Println("Deployment " + deploymentName + " has been deleted")
		return c.String(404, "Deployment "+deploymentName+" doesn't exist")
	}
}

func handleDeleteAutoscaler(c echo.Context) error {
	autoscalerName := c.Param("autoscalerName")

	if delete_api.DeleteAutoscaler(cli, autoscalerName) == true {
		fmt.Println("Autoscaler " + autoscalerName + " has been deleted")
		return c.String(200, "Autoscaler "+autoscalerName+" has been deleted")
	} else {
		fmt.Println("Autoscaler " + autoscalerName + " has been deleted")
		return c.String(404, "Autoscaler "+autoscalerName+" doesn't exist")
	}
}

func handleDeleteFuncPodInstance(c echo.Context) error {
	podInstanceID := c.Param("podName")
	flag, service := function_api.DeleteFuncPodInstance(cli, podInstanceID)
	if flag == true {
		fmt.Println("PodInstance of " + podInstanceID + " has been deleted")
		nodeList := get_api.GetAllNode(cli)
		for _, node := range nodeList {
			fmt.Println("service.ClusterIP is :   ", service.ClusterIP)
			letProxyDeleteCIRule(service.ClusterIP, node)
			letProxyCreateCIRule(service, node)
		}
		return c.String(200, "PodInstance of "+podInstanceID+" has been deleted")
	} else {
		return c.String(404, "PodInstance of "+podInstanceID+" doesn't exist")
	}
}

func handleDeleteFunction(c echo.Context) error {
	funcName := c.Param("funcName")
	flag, clusterIP := function_api.DeleteFunction(cli, funcName)
	if flag == true {
		fmt.Println("Function " + funcName + " has been deleted")
		nodeList := get_api.GetAllNode(cli)
		for _, node := range nodeList {
			fmt.Println("clusterIP is :   ", clusterIP)
			letProxyDeleteCIRule(clusterIP, node)
		}
		return c.String(200, "Function "+funcName+" has been deleted")
	} else {
		return c.String(404, "Function "+funcName+" doesn't exist")
	}
}

func handleDeletePodInstance(c echo.Context) error {
	podInstanceID := c.Param("podInstanceID")
	flag, service := function_api.FuncDeletePodInstance(cli, podInstanceID)
	if flag == true {
		fmt.Println("PodInstance of " + podInstanceID + " has been deleted")
		nodeList := get_api.GetAllNode(cli)
		for _, node := range nodeList {
			fmt.Println("service.ClusterIP is :   ", service.ClusterIP)
			letProxyDeleteCIRule(service.ClusterIP, node)
			letProxyCreateCIRule(service, node)
		}
		return c.String(200, "PodInstance of "+podInstanceID+" has been deleted")
	} else {
		return c.String(404, "PodInstance of "+podInstanceID+" doesn't exist")
	}
}

func handelDeleteStateMachine(c echo.Context) error {
	stateMachineName := c.Param("stateMachineName")

	if delete_api.DeleteStateMachine(cli, stateMachineName) == true {
		fmt.Println("StateMachine " + stateMachineName + " has been deleted")
		return c.String(200, "StateMachine "+stateMachineName+" has been deleted")
	} else {
		fmt.Println("StateMachine " + stateMachineName + " has been deleted")
		return c.String(404, "StateMachine "+stateMachineName+" doesn't exist")
	}
}

func letProxyDeleteCIRule(clusterIP string, node def.Node) {
	// 更新所有node的kube-proxy
	target := node.NodeIP.String() + ":" + strconv.Itoa(node.ProxyPort)
	fmt.Println("target is " + target)

	// 创建携程，并发执行
	//go func(target string) {
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
	//}(target)
}

func letProxyDeleteNPRule(clusterIP string, node def.Node) {
	// 更新所有node的kube-proxy
	target := node.NodeIP.String() + ":" + strconv.Itoa(node.ProxyPort)
	fmt.Println("target is " + target)

	// 创建携程，并发执行
	//go func(target string) {
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
	//}(target)
}

func handleGetAllNode(c echo.Context) error {
	nodeList := get_api.GetAllNodeInfo(cli)
	fmt.Println("nodeList:  ", nodeList)

	return c.JSON(200, nodeList)
}

func handleGetPod(c echo.Context) error {
	podInstanceName := c.Param("podInstanceName")
	podInstance, flag := get_api.GetPodInstance(cli, podInstanceName)
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
	fmt.Println("handleGetAllService")
	serviceList, flag := get_api.GetAllService(cli)

	if flag == false {
		return c.JSON(404, serviceList)
	}

	return c.JSON(200, serviceList)
}

func handleGetDeployment(c echo.Context) error {
	deploymentName := c.Param("deploymentName")
	deploymentDetail, flag := get_api.GetDeployment(cli, deploymentName)
	fmt.Println(deploymentDetail)

	if flag == false {
		return c.JSON(404, deploymentDetail)
	}

	return c.JSON(200, deploymentDetail)
}

func handleGetAllDeployment(c echo.Context) error {
	deploymentBriefList, flag := get_api.GetAllDeployment(cli)

	if flag == false {
		return c.JSON(404, deploymentBriefList)
	}

	return c.JSON(200, deploymentBriefList)
}

func handleGetAutoscaler(c echo.Context) error {
	autoscalerName := c.Param("autoscalerName")
	autoscalerDetail, flag := get_api.GetAutoscaler(cli, autoscalerName)
	fmt.Println(autoscalerDetail)

	if flag == false {
		return c.JSON(404, autoscalerDetail)
	}

	return c.JSON(200, autoscalerDetail)
}

func handleGetAllAutoscaler(c echo.Context) error {
	autoscalerBriefList, flag := get_api.GetAllAutoscaler(cli)

	if flag == false {
		return c.JSON(404, autoscalerBriefList)
	}

	return c.JSON(200, autoscalerBriefList)
}

func handleGetDNS(c echo.Context) error {
	dnsName := c.Param("dnsName")
	dnsDetail, flag := get_api.GetDNS(cli, dnsName)
	fmt.Println(dnsDetail)

	if flag == false {
		return c.JSON(404, dnsDetail)
	}

	return c.JSON(200, dnsDetail)
}

func handleGetAllDNS(c echo.Context) error {
	dnsDetailList, flag := get_api.GetAllDNS(cli)

	if flag == false {
		return c.JSON(404, dnsDetailList)
	}

	return c.JSON(200, dnsDetailList)
}

func handleGetFunction(c echo.Context) error {
	functionName := c.Param("functionName")
	function, flag := function_api.GetFunction(cli, functionName)
	fmt.Println(function)
	if flag == false {
		return c.JSON(404, function)
	}
	return c.JSON(200, function)
}

func handleGetAllFunction(c echo.Context) error {
	functionList, flag := function_api.GetAllFunction(cli)
	if flag == false {
		return c.JSON(404, functionList)
	}
	return c.JSON(200, functionList)
}

func handleGetGPUJob(c echo.Context) error {
	gpuJobName := c.Param("gpuJobName")
	gpuJobGet, flag := get_api.GetGPUJob(cli, gpuJobName)
	fmt.Println(gpuJobGet)
	if flag == false {
		return c.JSON(404, gpuJobGet)
	}
	return c.JSON(200, gpuJobGet)
}

func handleGetAllGPUJob(c echo.Context) error {
	gpuJobList, flag := get_api.GetAllGPUJob(cli)
	if flag == false {
		return c.JSON(404, gpuJobList)
	}
	return c.JSON(200, gpuJobList)
}

func handleGetStateMachine(c echo.Context) error {
	stateMachineName := c.Param("stateMachineName")
	stateMachine, flag := get_api.GetStateMachine(cli, stateMachineName)
	fmt.Println(stateMachine)
	if flag == false {
		return c.JSON(404, stateMachine)
	}
	return c.JSON(200, stateMachine)
}

func handleGetAllStateMachine(c echo.Context) error {
	stateMachineList, flag := get_api.GetAllStateMachine(cli)
	if flag == false {
		return c.JSON(404, stateMachineList)
	}
	return c.JSON(200, stateMachineList)
}

func NodeChecker() {
	cron2 := cron.New()
	err := cron2.AddFunc("*/180 * * * * *", CheckNode)
	if err != nil {
		fmt.Println(err)
	}

	cron2.Start()
	defer cron2.Stop()
	for {
		if ShouldStop {
			break
		}
	}
}

func CheckNode() {
	for nodeID, time_ := range HeartBeatMap {
		if time.Now().Sub(time_).Minutes() > 10 {
			node := apiserver_utils.GetNodeByID(cli, nodeID)
			node.Status = def.NotReady
			apiserver_utils.PersistNode(cli, node)
		}
	}
}
