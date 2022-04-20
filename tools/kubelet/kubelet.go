package kubelet

import (
	"fmt"
	"github.com/YOUR-USER-OR-ORG-NAME/YOUR-REPO-NAME/tools/def"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"os"
)

func main() {
	node := def.Node{}
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
	e.GET("/creatAndStartPod", creatAndStartPod)
	e.GET("/stopPod", stopPod)
	e.GET("/removePod", removePod)
	e.GET("/stopAndRemovePod", stopAndRemovePod)
	e.GET("/stopAll", stopAll)
	e.GET("/restartPod", restartPod)

	/*
		Goroutines
	*/
	go sendHeartbeat()

	e.Logger.Fatal(e.Start(":80"))
}

func stopAll(c echo.Context) error {
	return nil
}
