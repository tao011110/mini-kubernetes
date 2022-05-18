package def

import (
	"fmt"
	"github.com/jakehl/goid"
)

const (
	CadvisorPort                    = 8080
	EtcdPort                        = 2379
	ProxyPort                       = 3000
	GatewayImage                    = "hejingkai/zuul"
	GatewayRoutesConfigPathInImage  = `/home/zuul/src/main/resources/application.yaml`
	GatewayPackageAndRunScriptPath  = `/package_and_start.sh`
	NodeListName                    = `all_nodes_name`
	PodInstanceListID               = `pod_instance_list_id`
	DeploymentListName              = `deployment_list_name`
	PodListName                     = `pod_list_name`
	SchedulerPort                   = 9200
	ControllerPort                  = 8081
	NodeUndefined                   = -1
	HorizontalPodAutoscalerListName = `parsed_horizontal_pod_autoscaler_list_name`
	MasterIP                        = "192.168.1.7"
	MasterPort                      = "8000"
	EtcdDir                         = "/home/etcd-v3.2.13-linux-amd64"
	RgistryAddr                     = "registry.cn-hangzhou.aliyuncs.com/taoyucheng/mink8s:"
	RgistryUsername                 = "taoyucheng"
	RgistryPassword                 = "Tyc20010925tyc"
	TemplateImage                   = `hejingkai/python_serverless_template`
	PyHandlerPath                   = `/home/functionalTemplate/handler.py`
	RequirementsPath                = `/requirements.txt`
	PreparePath                     = `/prepare.sh`
	StartPath                       = `/start.sh`
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

func GenerateKeyOfPodInstanceReplicas(podInstanceName string) string {
	return GetKeyOfPodInstance(podInstanceName) + goid.NewV4UUID().String()
}

func GetKeyOfPodInstance(podInstanceName string) string {
	return fmt.Sprintf("/podInstance/%s-", podInstanceName)
}

func GetKeyOfDeployment(deploymentName string) string {
	return fmt.Sprintf("/deployment/%s", deploymentName)
}

func GetKeyOfAutoscaler(autoscalerName string) string {
	return fmt.Sprintf("/autoscaler/%s", autoscalerName)
}
