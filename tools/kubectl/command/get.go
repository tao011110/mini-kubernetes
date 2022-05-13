package command

import (
	"fmt"

	"mini-kubernetes/tools/def"
	"mini-kubernetes/tools/httpget"
	"mini-kubernetes/tools/master"

	"github.com/urfave/cli"
)

func NewGetCommand() cli.Command {
	return cli.Command{
		Name:  "get",
		Usage: "Get Pod state",
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
		err, status := httpget.Get("http://" + master.IP + ":" + master.Port + "/get/all/podStatus").
			ContentType("application/json").
			GetJson(&response).
			Execute()
		if err != nil {
			fmt.Println("[Fault] " + err.Error())
		}
		fmt.Printf("get_all_pod status is %s\n", status)
		if status == "200" {
			fmt.Println("All pods' brief information is as follows")
			for _, podInstanceBrief := range response {
				fmt.Printf("%v\n", podInstanceBrief)
			}
		} else {
			fmt.Printf("No pod exists\n")
		}
	} else if ty == "pod" && c.String("output") == "wide" {
		// kubectl get pod -o wide 查看全部Pod的状态
		response := make([]string, 0)
		err, status := httpget.Get("http://" + master.IP + ":" + master.Port + "/get/all/pod").
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
	}

}
