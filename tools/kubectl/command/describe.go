package command

import (
	"fmt"
	"mini-kubernetes/tools/def"
	"mini-kubernetes/tools/httpget"
	"mini-kubernetes/tools/util"

	"github.com/urfave/cli"
)

func NewDescribeCommand() cli.Command {
	return cli.Command{
		Name:  "create",
		Usage: "Create Pod according to xxx.yaml",
		Flags: []cli.Flag{},
		Action: func(c *cli.Context) error {
			describeFunc(c)
			return nil
		},
	}
}

func describeFunc(c *cli.Context) {
	if len(c.Args()) == 0 {
		fmt.Println("You need to specify pod name")
		return
	} else if (c.Args()[0] != "pod" && c.Args()[0] != "service") && c.Args()[0] != "dns" || len(c.Args()) < 2 {
		fmt.Println("Available command is 'kubectl describe pod/service.pod xxx'")
		return
	}

	if c.Args()[0] == "pod" {
		// kubectl describe pod podName
		podName := c.Args()[1]
		response := def.Pod{}
		err, status := httpget.Get("http://" + util.GetLocalIP().String() + ":" + fmt.Sprintf("%d", def.MasterPort) + "/get_pod/" + podName).
			ContentType("application/json").
			GetJson(&response).
			Execute()
		if err != nil {
			fmt.Println("[Fault] " + err.Error())
		}
		fmt.Printf("get_pod status is %s\n", status)
		if status == "200" {
			fmt.Printf("get pod_ %s successfully and the response is:\n %v\n", podName, response)
		} else {
			fmt.Printf("pod_ %s doesn't exist\n", podName)
		}
	} else if c.Args()[0] == "service" {
		// kubectl describe service serviceName
		// 用来获取特定名称的 service，需要发送给apiserver的参数为 serviceName(string)
		serviceName := c.Args()[1]
		response := def.Service{}
		err, status := httpget.Get("http://" + util.GetLocalIP().String() + ":" + fmt.Sprintf("%d", def.MasterPort) + "/get/service/" + serviceName).
			ContentType("application/json").
			GetJson(&response).
			Execute()
		if err != nil {
			fmt.Println("[Fault] " + err.Error())
		}
		fmt.Printf("get_service status is %s\n", status)
		if status == "200" {
			fmt.Printf("get service %s successfully and the response is: %v\n", serviceName, response)
		} else {
			fmt.Printf("service %s doesn't exist\n", serviceName)
		}
	} else if c.Args()[0] == "dns" {
		// kubectl describe dns dnsName
		// 用来获取特定名称的 dns，需要发送给apiserver的参数为 dnsName(string)
		//http调用返回的json需解析转为def.DNSDetail类型，
		dnsName := c.Args()[1]
		response := def.DNSDetail{}
		err, status := httpget.Get("http://" + util.GetLocalIP().String() + ":" + fmt.Sprintf("%d", def.MasterPort) + "/get/dns/" + dnsName).
			ContentType("application/json").
			GetJson(&response).
			Execute()
		if err != nil {
			fmt.Println("[Fault] " + err.Error())
		}
		fmt.Printf("get_dns status is %s\n", status)
		if status == "200" {
			fmt.Printf("get dns %s successfully and the response is: %v\n", dnsName, response)
		} else {
			fmt.Printf("dns %s doesn't exist\n", dnsName)
		}
	}

}
