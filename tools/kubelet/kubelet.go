package kubelet

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"mini-kubernetes/tools/def"
	"mini-kubernetes/tools/docker"
	"mini-kubernetes/tools/etcd"
	"mini-kubernetes/tools/kubeproxy"
	"mini-kubernetes/tools/resource"
	"os"
)

var node = def.Node{}

func main() {
	parseArgs(&node.NodeName, &node.MasterIpAndPort, &node.LocalPort)
	node.NodeIP = getLocalIP()
	node.ProxyPort = kubeproxy.ProxyPort
	if node.NodeIP == nil {
		fmt.Println("get local ip error")
		os.Exit(0)
	}
	err := registerToMaster(&node)
	if err == nil {
		fmt.Println("network error, cannot register to master")
		os.Exit(0)
	}
	docker.CreateNetBridge(node.CniIP.String())

	/*
		creat echo instance
	*/
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	etcdClient, err := etcd.Start("", def.EtcdPort)
	if err != nil {
		e.Logger.Error("Start etcd error!")
		os.Exit(0)
	}
	node.EtcdClient = etcdClient
	cadvisorClient, err := resource.StartCadvisor()
	if err != nil {
		e.Logger.Error("Start cadvisor error!")
		os.Exit(0)
	}
	node.CadvisorClient = cadvisorClient

	go EtcdWatcher()
	go ResourceMonitoring()
	go ContainerCheck()

	e.Logger.Fatal(e.Start(":80"))

}
