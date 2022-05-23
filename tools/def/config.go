package def

import (
	"fmt"
	"github.com/jakehl/goid"
)

const (
	NodeUndefined = -1
	EtcdDir       = "/home/etcd"
)

/********** HTTP ports **********/
const (
	CadvisorPort   = 8080
	EtcdPort       = 2379
	ProxyPort      = 3000
	SchedulerPort  = 9200
	ControllerPort = 8081
	ActiverPort    = 3306
	MasterPort     = 8000
)

/********** Image(gateway and functional) **********/
const (
	//for image upload
	RgistryAddr     = "registry.cn-hangzhou.aliyuncs.com/taoyucheng/mink8s:"
	RgistryUsername = "taoyucheng"
	RgistryPassword = "Tyc20010925tyc"

	//for gateway image
	GatewayImage                   = "hejingkai/zuul"
	GatewayRoutesConfigPathInImage = `/home/zuul/src/main/resources/application.yaml`
	GatewayPackageCmd              = `./package.sh`
	GatewayStartCmd                = `./start.sh`

	//for py functional image
	PyFunctionTemplateImage   = `hejingkai/python_serverless_template`
	PyHandlerParentDirPath    = `/home/functionalTemplate/`
	PyHandlerFileName         = `handler.py`
	RequirementsParentDirPath = `/`
	RequirementsFileName      = `requirements.txt`
	PyFunctionPrepareCmd      = `./prepare.sh`
	PyFunctionStartCmd        = `./start.sh`
	MaxBodySize               = 2048

	//for gpu job uploader image
	GPUJobUploaderImage                    = `hejingkai/gpu_job_uploader`
	GPUJobSourceCodeParentDirPath          = `/home/gpu/`
	GPUJobSourceCodeFileName               = `cuda_mat_add.cu`
	GPUJOBMakefileParentDirPath            = `/home/gpu/`
	GPUJOBMakefileFileName                 = `Makefile`
	GPUSlurmScriptParentDirPath            = `/home/gpu/`
	GPUSlurmScriptFileName                 = `job.slurm`
	GPUApiServerIpAndPortFileParentDirPath = `/home/result/`
	GPUApiServerIpAndPortFileFileName      = `apiserver_ip_and_port`
	GPUJobNameParentDirName                = `/home/`
	GPUJobNameFileName                     = `job_name`
	GPUJobUploaderRunCmd                   = `./home/run.sh`

	TemplateCmdFilePath = "/home/temp_cmd.sh"
	//TemplateFileDir     = `/home/temp/`
	//ONLY FOR TEST !!!
	TemplateFileDir = `/home/hjk/Desktop/Go/mini-kubernetes/test/slurmGenrator_test/`
)

/********** ETCD key **********/
const (
	NodeListID                      = `all_nodes_id`
	PodInstanceListID               = `pod_instance_list_id`
	DeploymentListName              = `deployment_list_name`
	FunctionNameListKey             = `function_name_list`
	StateMachineNameListKey         = `state_machine_name_list_key`
	HorizontalPodAutoscalerListName = `parsed_horizontal_pod_autoscaler_list_name`
)

func GetKeyOfPodReplicasNameListByPodName(podName string) string {
	return fmt.Sprintf("%s_replicas_name_list", podName)
}

func PodInstanceListKeyOfNode(node *Node) string {
	return fmt.Sprintf("node%d_pod_instances", node.NodeID)
}

func PodInstanceListKeyOfNodeID(nodeID int) string {
	return fmt.Sprintf("node%d_pod_instances", nodeID)
}

func KeyNodeResourceUsage(nodeID int) string {
	return fmt.Sprintf("%d_resource_usage", nodeID)
}

func GetKeyOfResourceUsageByPodInstanceID(instanceID string) string {
	return fmt.Sprintf("%s_resource_usage", instanceID)
}

func GetKeyOfPod(podName string) string {
	return fmt.Sprintf("/pod/%s", podName)
}

func GetKeyOfService(serviceName string) string {
	return fmt.Sprintf("/service/%s", serviceName)
}

func GetKeyOfFunction(name string) string {
	return fmt.Sprintf("/function/%s", name)
}

func GetKeyOfStateMachine(name string) string {
	return fmt.Sprintf("/state_machine/%s", name)
}

func GenerateKeyOfPodInstanceReplicas(podInstanceName string) string {
	return GetKeyOfPodInstance(podInstanceName) + "-" + goid.NewV4UUID().String()
}

func GetKeyOfPodInstance(podInstanceName string) string {
	return fmt.Sprintf("/podInstance/%s", podInstanceName)
}

func GetKeyOfDeployment(deploymentName string) string {
	return fmt.Sprintf("/deployment/%s", deploymentName)
}

func GetKeyOfAutoscaler(autoscalerName string) string {
	return fmt.Sprintf("/autoscaler/%s", autoscalerName)
}

func GetPodNameOfDeployment(deploymentName string) string {
	return fmt.Sprintf("%s-pod", deploymentName)
}

func GetPodNameOfAutoscaler(autoscalerName string) string {
	return fmt.Sprintf("%s-pod", autoscalerName)
}
