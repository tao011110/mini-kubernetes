package def

import "fmt"

const (
	CadvisorPort                    = 8080
	EtcdPort                        = 2379
	ProxyPort                       = 3000
	GatewayImage                    = "hejingkai/zuul"
	GatewayRoutesConfigPathInImage  = `/home/zuul/src/main/resources/application.yaml`
	GatewayPackageAndRunScriptPath  = `/package_and_start.sh`
	NodeListName                    = `all_nodes_name`
	PodInstanceListName             = `pod_instance_list_name`
	DeploymentListName              = `deployment_list_name`
	PodListName                     = `pod_list_name`
	SchedulerPort                   = 9200
	ControllerPort                  = 8081
	NodeUndefined                   = -1
	HorizontalPodAutoscalerListName = `parsed_horizontal_pod_autoscaler_list_name`
	MasterIP                        = "192.168.1.7"
	MasterPort                      = "8000"
	EtcdDir                         = "/home/etcd-v3.2.13-linux-amd64"
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
