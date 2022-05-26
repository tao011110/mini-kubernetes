package command

import (
	"fmt"
	"mini-kubernetes/tools/util"

	"mini-kubernetes/tools/def"
	"mini-kubernetes/tools/httpget"

	"github.com/urfave/cli"
)

func NewGetCommand() cli.Command {
	return cli.Command{
		Name:  "get",
		Usage: "Get resources state",
		Flags: []cli.Flag{
			cli.StringFlag{Name: "output, o", Value: "", Usage: "Output format"},
		},
		Action: func(c *cli.Context) error {
			getFunc(c)
			return nil
		},
	}
}

func getFunc(c *cli.Context) {

	if len(c.Args()) == 0 {
		fmt.Println("You need to specify get what")
		return
	}

	ty := c.Args()[0]
	if ty == "pods" {
		// kubectl get pods 查看全部Pod的概要状态
		response := make([]def.PodInstanceBrief, 0)
		err, status := httpget.Get("http://" + util.GetLocalIP().String() + ":" + fmt.Sprintf("%d", def.MasterPort) + "/get/all/podStatus").
			ContentType("application/json").
			GetJson(&response).
			Execute()
		if err != nil {
			fmt.Println("[Fault] " + err.Error())
		}
		fmt.Printf("get_all_pod status is %s\n", status)
		if status == "200" {
			fmt.Println("All pods' brief information is as follows")
			fmt.Println("NAME   READY   STATUS   RESTARTS   AGE")
			for _, podInstanceBrief := range response {
				fmt.Printf("%v\n", podInstanceBrief)
			}
		} else {
			fmt.Printf("No pod exists\n")
		}
	} else if ty == "pods" && c.String("output") == "wide" {
		// kubectl get pods -o wide 查看全部Pod的状态
		response := make([]string, 0)
		err, status := httpget.Get("http://" + util.GetLocalIP().String() + ":" + fmt.Sprintf("%d", def.MasterPort) + "/get/all/pod").
			ContentType("application/json").
			GetJson(&response).
			Execute()
		if err != nil {
			fmt.Println("[Fault] " + err.Error())
		}
		fmt.Printf("get_all_pod status is %s\n", status)
		if status == "200" {
			fmt.Println("All pods are as follows")
			for _, podInstance := range response {
				fmt.Printf("%v\n", podInstance)
			}
		} else {
			fmt.Printf("No pod exists\n")
		}
	} else if ty == "services" {
		// kubectl get services
		// 用来获取所有的 service
		response := make([]def.Service, 0)
		err, status := httpget.Get("http://" + util.GetLocalIP().String() + ":" + fmt.Sprintf("%d", def.MasterPort) + "/get/all/service").
			ContentType("application/json").
			GetJson(&response).
			Execute()
		if err != nil {
			fmt.Println("[Fault] " + err.Error())
		}
		fmt.Printf("get_all_service status is %s\n", status)
		if status == "200" {
			fmt.Println("All services' information is as follows")
			for _, service := range response {
				fmt.Printf("%v\n", service)
			}
		} else {
			fmt.Printf("No service exists\n")
		}
	} else if ty == "dns" {
		// kubectl get dns
		// 用来获取所有的 dns
		response := make([]def.DNSDetail, 0)
		err, status := httpget.Get("http://" + util.GetLocalIP().String() + ":" + fmt.Sprintf("%d", def.MasterPort) + "/get/all/dns").
			ContentType("application/json").
			GetJson(&response).
			Execute()
		if err != nil {
			fmt.Println("[Fault] " + err.Error())
		}
		fmt.Printf("get_all_dns status is %s\n", status)
		if status == "200" {
			fmt.Println("All dns' information is as follows")
			for _, service := range response {
				fmt.Printf("%v\n", service)
			}
		} else {
			fmt.Printf("No dns exists\n")
		}
	} else if ty == "deployment" {
		// kubectl get deployment 用来获取所有的 deployment
		// DeploymentBrief提供了显示需要的全部信息
		response := make([]def.DeploymentBrief, 0)
		err, status := httpget.Get("http://" + util.GetLocalIP().String() + ":" + fmt.Sprintf("%d", def.MasterPort) + "/get/all/deployment").
			ContentType("application/json").
			GetJson(&response).
			Execute()
		if err != nil {
			fmt.Println("[Fault] " + err.Error())
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
	} else if ty == "autoscaler" {
		// 用来获取所有的 autoscaler
		// AutoscalerBrief提供了 的 kubelet get autoscaler 显示的部分信息
		response := make([]def.AutoscalerBrief, 0)
		err, status := httpget.Get("http://" + util.GetLocalIP().String() + ":" + fmt.Sprintf("%d", def.MasterPort) + "/get/all/autoscaler").
			ContentType("application/json").
			GetJson(&response).
			Execute()
		if err != nil {
			fmt.Println("[Fault] " + err.Error())
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

}
