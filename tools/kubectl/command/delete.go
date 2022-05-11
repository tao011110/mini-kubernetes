package command

import (
	"fmt"

	"mini-kubernetes/tools/httpget"
	"mini-kubernetes/tools/master"
	"mini-kubernetes/tools/yaml"

	"github.com/urfave/cli"
)

func NewDeleteCommand() cli.Command {
	return cli.Command{
		Name:  "delete",
		Usage: "Delete Pod according to xxx.yaml or name",
		Flags: []cli.Flag{
			cli.StringFlag{Name: "file, f", Value: "", Usage: "File path to the config"},
		},
		Action: func(c *cli.Context) error {
			deleteFunc(c)
			return nil
		},
	}
}

func deleteFunc(c *cli.Context) {

	dir := c.String("file")
	var podName string
	if dir != "" {
		// 根据yaml文件确定pod名称
		fmt.Printf("Using dir: %s\n", dir)
		pod_, _ := yaml.ReadPodYamlConfig(dir)
		podName = pod_.Metadata.Name
	} else {
		if len(c.Args()) == 0 {
			fmt.Println("You need to specify pod name")
			return
		}
		podName = c.Args()[0]
		fmt.Printf("Podname: %s\n", podName)
	}

	// delete_pod
	// 需要发送给apiserver的参数为 podName string
	response := ""
	err, status := httpget.DELETE("http://" + master.IP + ":" + master.Port + "/delete_pod/" + podName).
		ContentType("application/json").
		GetString(&response).
		Execute()
	if err != nil {
		fmt.Println("[Fault] " + err.Error())
	} else {
		fmt.Printf("get_pod status is %s\n", status)
		if status == "200" {
			fmt.Printf("delete pod_ %s successfully and the response is: %v\n", podName, response)
		} else {
			fmt.Printf("pod_ %s doesn't exist\n", podName)
		}
	}
}
