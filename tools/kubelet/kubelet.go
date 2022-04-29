package kubelet

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"mini-kubernetes/tools/def"
	"mini-kubernetes/tools/etcd"
	"mini-kubernetes/tools/resource"
	"os"
)

/*
	shared between all routines
*/
var node = def.Node{}

func main() {
	parseArgs(&node.NodeName, &node.MasterIpAndPort, &node.LocalPort)
	node.NodeIP = getLocalIP()
	if node.NodeIP == nil {
		fmt.Println("get local ip error")
		os.Exit(0)
	}
	err := registerToMaster(&node)
	if err == nil {
		fmt.Println("network error, cannot register to master")
		os.Exit(0)
	}

	/*
		creat echo instance
	*/
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	/*
		register handlers
	*/
	etcdClient, err := etcd.Start("", def.EtcdPort)
	if err != nil {
		e.Logger.Error("Start etcd error!")
		os.Exit(0)
	}
	cadvisorClient, err := resource.StartCadvisor()
	if err != nil {
		e.Logger.Error("Start cadvisor error!")
		os.Exit(0)
	}

	go EtcdWatcher(etcdClient, node.NodeID, e)

	e.Logger.Fatal(e.Start(":80"))

}

func stopAll(c echo.Context) error {
	return nil
}
