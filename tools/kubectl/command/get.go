package command

import (
	"fmt"
	"mini-kubernetes/tools/util"
	"strconv"
	"strings"
	"time"

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
		// fmt.Printf("get_all_pod status is %s\n", status)
		if status == "200" {
			fmt.Println("All pods' brief information is as follows")
			max := 12
			fmt.Printf("%-"+strconv.Itoa(max)+"s %-"+strconv.Itoa(max)+"s %-"+strconv.Itoa(max)+"s %-"+strconv.Itoa(max)+"s %-"+strconv.Itoa(max)+"s\n",
				"NAME",
				"READY",
				"STATUS",
				"RESTARTS",
				"AGE")
			for _, podInstanceBrief := range response {
				fmt.Printf("%-"+strconv.Itoa(max)+"s %-"+strconv.Itoa(max)+"s %-"+strconv.Itoa(max)+"s %-"+strconv.Itoa(max)+"s %-"+strconv.Itoa(max)+"s\n",
					podInstanceBrief.Name,
					podInstanceBrief.Ready,
					strconv.Itoa(int(podInstanceBrief.Status)),
					strconv.Itoa(int(podInstanceBrief.Restarts)),
					podInstanceBrief.Age)
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
		// fmt.Printf("get_all_pod status is %s\n", status)
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
		// fmt.Printf("get_all_service status is %s\n", status)
		if status == "200" {
			fmt.Println("All services' information is as follows")
			max := 12
			ip := 15
			fmt.Printf("%-"+strconv.Itoa(max)+"s %-"+strconv.Itoa(max)+"s %-"+strconv.Itoa(ip)+"s %-"+
				strconv.Itoa(ip)+"s %-"+strconv.Itoa(max)+"s %-"+strconv.Itoa(max)+"s\n",
				"NAME",
				"TYPE",
				"CLUSTER-IP",
				"EXTERNAL-IP",
				"PORT(S)",
				"AGE")
			for _, service := range response {
				t := time.Now()                                  // 用于获取当前时间
				var Age time.Duration = t.Sub(service.StartTime) //进行计算，得到AGE
				for i := range service.PortsBindings {
					fmt.Printf("%-"+strconv.Itoa(max)+"s %-"+strconv.Itoa(max)+"s %-"+strconv.Itoa(ip)+"s %-"+
						strconv.Itoa(ip)+"s %-"+strconv.Itoa(max)+"s %-"+strconv.Itoa(max)+"s\n",
						service.Name,
						service.Type,
						service.ClusterIP,
						"<none>",
						strconv.Itoa(int(service.PortsBindings[i].Ports.Port))+"/"+strings.ToUpper(service.PortsBindings[i].Ports.Protocol),
						Age)
				}
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
		// fmt.Printf("get_all_dns status is %s\n", status)
		if status == "200" {
			fmt.Println("All dns' information is as follows")
			max := 15
			fmt.Printf("%-"+strconv.Itoa(max)+"s %-"+strconv.Itoa(max)+"s %-"+strconv.Itoa(max)+"s %-"+
				strconv.Itoa(max)+"s %-"+strconv.Itoa(max)+"s\n",
				"NAME",
				"HOST",
				"PATH",
				"SERVICE-NAME",
				"PORT")
			for _, dns := range response {
				for i := range dns.Paths {
					fmt.Printf("%-"+strconv.Itoa(max)+"s %-"+strconv.Itoa(max)+"s %-"+strconv.Itoa(max)+"s %-"+
						strconv.Itoa(max)+"s %-"+strconv.Itoa(max)+"s\n",
						dns.Name,
						dns.Host,
						dns.Paths[i].Path,
						dns.Paths[i].Service.Metadata.Name,
						strconv.Itoa(int(dns.Paths[i].Port)))
				}
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
		// fmt.Printf("get_all_deployment status is %s\n", status)
		if status == "200" {
			fmt.Println("All deployments' information is as follows")
			max := 15
			num := 10
			fmt.Printf("%-"+strconv.Itoa(max)+"s %-"+strconv.Itoa(max)+"s %-"+strconv.Itoa(num)+"s %-"+
				strconv.Itoa(num)+"s %-"+strconv.Itoa(max)+"s\n",
				"NAME",
				"READY",
				"UpToDate",
				"AVAILABLE",
				"AGE")
			for _, dep := range response {
				fmt.Printf("%-"+strconv.Itoa(max)+"s %-"+strconv.Itoa(max)+"s %-"+strconv.Itoa(num)+"s %-"+
					strconv.Itoa(num)+"s %-"+strconv.Itoa(max)+"s\n",
					dep.Name,
					dep.Ready,
					strconv.Itoa(dep.UpToDate),
					strconv.Itoa(dep.Available),
					dep.Age)
			}
		} else {
			fmt.Printf("No deployment exists\n")
		}
	} else if ty == "autoscaler" {
		// kubectl get autoscaler 用来获取所有的 autoscaler
		// AutoscalerBrief提供了 的 kubelet get autoscaler 显示的部分信息
		response := make([]def.AutoscalerBrief, 0)
		err, status := httpget.Get("http://" + util.GetLocalIP().String() + ":" + fmt.Sprintf("%d", def.MasterPort) + "/get/all/autoscaler").
			ContentType("application/json").
			GetJson(&response).
			Execute()
		if err != nil {
			fmt.Println("[Fault] " + err.Error())
		}
		// fmt.Printf("get_all_autoscaler status is %s\n", status)
		if status == "200" {
			fmt.Println("All autoscalers' information is as follows")
			max := 15
			num := 10
			fmt.Printf("%-"+strconv.Itoa(max)+"s %-"+strconv.Itoa(num)+"s %-"+strconv.Itoa(num)+"s %-"+
				strconv.Itoa(num)+"s %-"+strconv.Itoa(max)+"s\n",
				"NAME",
				"MIN-PODS",
				"MAX-PODS",
				"REPLICAS",
				"AGE")
			for _, auto := range response {
				fmt.Printf("%-"+strconv.Itoa(max)+"s %-"+strconv.Itoa(num)+"s %-"+strconv.Itoa(num)+"s %-"+
					strconv.Itoa(num)+"s %-"+strconv.Itoa(max)+"s\n",
					auto.Name,
					strconv.Itoa(auto.MinPods),
					strconv.Itoa(auto.MaxPods),
					strconv.Itoa(auto.Replicas),
					auto.Age)
			}
		} else {
			fmt.Printf("No autoscaler exists\n")
		}
	} else if ty == "gpujob" {
		// kubectl get gpujob 用来获取所有的 gpuJob
		response := make([]def.GPUJobDetail, 0)
		err, status := httpget.Get("http://" + util.GetLocalIP().String() + ":" + fmt.Sprintf("%d", def.MasterPort) + "/get/all/gpuJob").
			ContentType("application/json").
			GetJson(&response).
			Execute()
		if err != nil {
			fmt.Println("[Fault] " + err.Error())
		}
		// fmt.Printf("get_all_gpuJob status is %s\n", status)
		if status == "200" {
			fmt.Println("All gpuJobs' information is as follows")
			for _, gpuJobDetail := range response {
				fmt.Printf("%v\n", gpuJobDetail)
			}
		} else {
			fmt.Printf("No gpuJob exists\n")
		}
	} else if ty == "function" {
		// kubectl get function 用来获取所有的 function
		response := make([]def.Function, 0)
		err, status := httpget.Get("http://" + util.GetLocalIP().String() + ":" + fmt.Sprintf("%d", def.MasterPort) + "/get/all/function").
			ContentType("application/json").
			GetJson(&response).
			Execute()
		if err != nil {
			fmt.Println("[Fault] " + err.Error())
		}
		// fmt.Printf("get_all_function status is %s\n", status)
		if status == "200" {
			fmt.Println("All functions' information is as follows")
			for _, function := range response {
				fmt.Printf("%v\n", function)
			}
		} else {
			fmt.Printf("No function exists\n")
		}
	} else if ty == "statemachine" {
		// kubectl get statemachine 用来获取所有的 statemachine
		response := make([]def.StateMachine, 0)
		err, status := httpget.Get("http://" + util.GetLocalIP().String() + ":" + fmt.Sprintf("%d", def.MasterPort) + "/get/all/stateMachine").
			ContentType("application/json").
			GetJson(&response).
			Execute()
		if err != nil {
			fmt.Println("[Fault] " + err.Error())
		}
		// fmt.Printf("get_all_stateMachine status is %s\n", status)
		if status == "200" {
			fmt.Println("All stateMachines' information is as follows")
			for _, stateMachine := range response {
				fmt.Printf("%v\n", stateMachine)
			}
		} else {
			fmt.Printf("No stateMachine exists\n")
		}
	}
}
