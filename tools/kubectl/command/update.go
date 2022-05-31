package command

import (
	"fmt"
	"bytes"
	"encoding/json"
	"mini-kubernetes/tools/util"
	"mini-kubernetes/tools/def"
	"mini-kubernetes/tools/httpget"
	"mini-kubernetes/tools/yaml"

	"github.com/urfave/cli"
)

func NewUpdateCommand() cli.Command {
	return cli.Command{
		Name:  "update",
		Usage: "Update function",
		Flags: []cli.Flag{
			cli.StringFlag{Name: "soft, s", Value: "", Usage: "Soft update function"},
			cli.StringFlag{Name: "hard, h", Value: "", Usage: "Hard update function"},
		},
		Action: func(c *cli.Context) error {
			updateFunc(c)
			return nil
		},
	}
}

func updateFunc(c *cli.Context) {

	dir_s := c.String("soft")
	dir_h := c.String("hard")
	
	if dir_s != "" && dir_h == "" {
		// 用来软更新function，需要发送给apiserver的参数为 function (def.Function)
		fmt.Printf("Using dir: %s\n", dir_s)
		function, err := yaml.ReadFunctionConfig(dir_s)
		if function == nil || err != nil {
			fmt.Println("[Fault] " + err.Error())
			return
		}
		request := *function
		response := ""
		body, _ := json.Marshal(request)
		err, status := httpget.Put("http://" + util.GetLocalIP().String() + ":" + fmt.Sprintf("%d", def.MasterPort) + "/update/soft/function").
			ContentType("application/json").
			Body(bytes.NewReader(body)).
			GetString(&response).
			Execute()
		if err != nil {
			fmt.Println("[Fault] " + err.Error())
		}
		fmt.Printf("soft update function is %s and response is: %s\n", status, response)

	} else if dir_s == "" && dir_h != "" {
		// 用来硬更新function，需要发送给apiserver的参数为 function (def.Function)
		fmt.Printf("Using dir: %s\n", dir_h)
		function, err := yaml.ReadFunctionConfig(dir_h)
		if function == nil || err != nil {
			fmt.Println("[Fault] " + err.Error())
			return
		}
		request := *function
		response := ""
		body, _ := json.Marshal(request)
		err, status := httpget.Put("http://" + util.GetLocalIP().String() + ":" + fmt.Sprintf("%d", def.MasterPort) + "/update/hard/function").
			ContentType("application/json").
			Body(bytes.NewReader(body)).
			GetString(&response).
			Execute()
		if err != nil {
			fmt.Println("[Fault] " + err.Error())
		}
		fmt.Printf("hard update function is %s and response is: %s\n", status, response)
	}
}