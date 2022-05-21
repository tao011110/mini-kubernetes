package command

import (
	"fmt"
	"mini-kubernetes/tools/def"
	"mini-kubernetes/tools/httpget"

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
	} else if (c.Args()[0] != "pod" && c.Args()[0] != "service") || len(c.Args()) < 2 {
		fmt.Println("Available command is 'kubectl describe pod/service xxx'")
		return
	}

	if c.Args()[0] == "pod" {
		// kubectl describe pod podName
		podName := c.Args()[1]
		response := def.Pod{}
		err, status := httpget.Get("http://" + def.MasterIP + ":" + def.MasterPort + "/get_pod/" + podName).
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
	} else {
		// kubectl describe service serviceName
		// 用来获取特定名称的 service，需要发送给apiserver的参数为 serviceName(string)
		serviceName := c.Args()[1]
		response := def.Service{}
		err, status := httpget.Get("http://" + def.MasterIP + ":" + def.MasterPort + "/get/service/" + serviceName).
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
	}

}
