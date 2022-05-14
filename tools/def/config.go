package def

const (
	CadvisorPort                   = 8080
	EtcdPort                       = 2379
	GatewayImage                   = "hejingkai/zuul"
	GatewayRoutesConfigPathInImage = `/home/zuul/src/main/resources/application.yaml`
	GatewayPackageAndRunScriptPath = `/package_and_start.sh`
	NodeListName                   = `all_nodes_name`
	PodInstanceListName            = `pod_instance_list_name`
)
