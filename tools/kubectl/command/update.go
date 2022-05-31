package command

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mini-kubernetes/tools/def"
	"mini-kubernetes/tools/httpget"
	"mini-kubernetes/tools/util"
	"mini-kubernetes/tools/yaml"

	"github.com/urfave/cli"
)

func NewUpdateCommand() cli.Command {
	return cli.Command{
		Name:  "update",
		Usage: "Update function",
		Flags: []cli.Flag{
			cli.StringFlag{Name: "file, f", Value: "", Usage: "File path to the config"},
		},
		Action: func(c *cli.Context) error {
			updateFunc(c)
			return nil
		},
	}
}

func updateFunc(c *cli.Context) {

	dir := c.String("file")
	if dir == "" {
		wrong("You need to specify directory")
		return
	}

	if c.Args()[0] == "soft" {
		// 用来软更新function，需要发送给apiserver的参数为 function (def.Function)
		fmt.Printf("Using dir: %s\n", dir)
		function, err := yaml.ReadFunctionConfig(dir)
		if function == nil || err != nil {
			wrong(err.Error())
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
			wrong(err.Error())
		}
		fmt.Printf("soft update function is %s and response is: %s\n", status, response)

	} else if c.Args()[0] == "hard" {
		// 用来硬更新function，需要发送给apiserver的参数为 function (def.Function)
		fmt.Printf("Using dir: %s\n", dir)
		function, err := yaml.ReadFunctionConfig(dir)
		if function == nil || err != nil {
			wrong(err.Error())
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
			wrong(err.Error())
		}
		fmt.Printf("hard update function is %s and response is: %s\n", status, response)
	} else {
		wrong("Wrong update type")
	}
}
