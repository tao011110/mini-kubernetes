package command

import (
	"fmt"
	"github.com/urfave/cli"
	"mini-kubernetes/tools/def"
	"mini-kubernetes/tools/httpget"
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
	} else if c.Args()[0] != "pod" || len(c.Args()) < 2 {
		fmt.Println("Available command is 'kubectl describe pod xxx'")
		return
	}

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

}
