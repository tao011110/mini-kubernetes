package command

import (
	"bytes"
	"encoding/json"
	"fmt"

	"mini-kubernetes/tools/httpget"
	"mini-kubernetes/tools/master"
	"mini-kubernetes/tools/yaml"

	"github.com/urfave/cli"
)

func NewCreateCommand() cli.Command {
	return cli.Command{
		Name:  "create",
		Usage: "Create Pod according to xxx.yaml",
		Flags: []cli.Flag{
			cli.StringFlag{Name: "file, f", Value: "", Usage: "File path to the config"},
		},
		Action: func(c *cli.Context) error {
			createFunc(c)
			return nil
		},
	}
}

func createFunc(c *cli.Context) {
	// create_pod
	// 需要发送给apiserver的参数为 pod_ def.Pod
	dir := c.String("file")
	fmt.Printf("Using dir: %s\n", dir)
	pod_, _ := yaml.ReadPodYamlConfig(dir)

	request := *pod_
	response := ""
	body, _ := json.Marshal(request)
	err, status := httpget.Post("http://" + master.IP + ":" + master.Port + "/create_pod").
		ContentType("application/json").
		Body(bytes.NewReader(body)).
		GetString(&response).
		Execute()
	if err != nil {
		fmt.Println("[Fault] " + err.Error())
	} else {
		fmt.Printf("create_pod is %s and response is: %s\n", status, response)
	}
}
