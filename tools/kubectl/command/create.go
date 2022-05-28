package command

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mini-kubernetes/tools/def"
	"mini-kubernetes/tools/httpget"
	"mini-kubernetes/tools/util"
	"mini-kubernetes/tools/yaml"
	"os"

	"github.com/urfave/cli"
)

func NewCreateCommand() cli.Command {
	return cli.Command{
		Name:  "create",
		Usage: "Create resources according to xxx.yaml",
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

	dir := c.String("file")
	fmt.Printf("Using dir: %s\n", dir)
	last4 := dir[len(dir)-4:]
	if last4 == "json" {
		// 用来创建StateMachine，需要发送给apiserver的参数为 stateMachine (def.StateMachine)
		// 直接读取文件
		file, err := os.Open(dir)
		if err != nil {
			panic(err)
		}
		defer func(file *os.File) {
			_ = file.Close()
		}(file)
		body, _ := ioutil.ReadAll(file)

		response := ""
		err, status := httpget.Post("http://" + util.GetLocalIP().String() + ":" + fmt.Sprintf("%d", def.MasterPort) + "/create/stateMachine").
			ContentType("application/json").
			Body(bytes.NewReader(body)).
			GetString(&response).
			Execute()
		if err != nil {
			fmt.Println("[Fault] " + err.Error())
		}
		fmt.Printf("create_stateMachine is %s and response is: %s\n", status, response)

	} else if last4 == "yaml" {
		ty, _ := yaml.ReadType(dir)
		switch ty {
		case yaml.Pod_t:
			// create_pod
			// 需要发送给apiserver的参数为 pod_ def.Pod
			pod_, err := yaml.ReadPodYamlConfig(dir)
			if pod_ == nil {
				fmt.Println("[Fault] " + err.Error())
			}
			request := *pod_
			response := ""
			body, _ := json.Marshal(request)
			err, status := httpget.Post("http://" + util.GetLocalIP().String() + ":" + fmt.Sprintf("%d", def.MasterPort) + "/create_pod").
				ContentType("application/json").
				Body(bytes.NewReader(body)).
				GetString(&response).
				Execute()
			if err != nil {
				fmt.Println("[Fault] " + err.Error())
			} else {
				fmt.Printf("create_pod is %s and response is: %s\n", status, response)
			}
		case yaml.ClusterIP_t:
			// create_clusterIP_service
			serviceC_, err := yaml.ReadServiceClusterIPConfig(dir)
			if serviceC_ == nil {
				fmt.Println("[Fault] " + err.Error())
			}
			request := *serviceC_
			response := ""
			body, _ := json.Marshal(request)
			err, status := httpget.Post("http://" + util.GetLocalIP().String() + ":" + fmt.Sprintf("%d", def.MasterPort) + "/create/clusterIPService").
				ContentType("application/json").
				Body(bytes.NewReader(body)).
				GetString(&response).
				Execute()
			if err != nil {
				fmt.Println("[Fault] " + err.Error())
			} else {
				fmt.Printf("create_service is %s and response is: %s\n", status, response)
			}
		case yaml.Nodeport_t:
			// create_nodeport_service
			serviceN_, err := yaml.ReadServiceNodeportConfig(dir)
			if serviceN_ == nil {
				fmt.Println("[Fault] " + err.Error())
			}
			request := *serviceN_
			response := ""
			body, _ := json.Marshal(request)
			err, status := httpget.Post("http://" + util.GetLocalIP().String() + ":" + fmt.Sprintf("%d", def.MasterPort) + "/create/nodePortService").
				ContentType("application/json").
				Body(bytes.NewReader(body)).
				GetString(&response).
				Execute()
			if err != nil {
				fmt.Println("[Fault] " + err.Error())
			}
			fmt.Printf("create_service is %s and response is: %s\n", status, response)
		case yaml.Dns_t:
			// 用来创建DNS和Gateway, 需要发送给apiserver的参数为 dns def.DNS
			dns, _ := yaml.ReadDNSConfig(dir)
			request := *dns
			response := ""
			body, _ := json.Marshal(request)
			err, status := httpget.Post("http://" + util.GetLocalIP().String() + ":" + fmt.Sprintf("%d", def.MasterPort) + "/create/dns").
				ContentType("application/json").
				Body(bytes.NewReader(body)).
				GetString(&response).
				Execute()
			if err != nil {
				fmt.Println("[Fault] " + err.Error())
			}
			fmt.Printf("create_gateway is %s and response is: %s\n", status, response)
		case yaml.Deployment_t:
			// 用来创建Deployment，需要发送给apiserver的参数为 deployment (def.Deployment)
			deployment, _ := yaml.ReadDeploymentConfig(dir)
			request := *deployment
			response := ""
			body, _ := json.Marshal(request)
			err, status := httpget.Post("http://" + util.GetLocalIP().String() + ":" + fmt.Sprintf("%d", def.MasterPort) + "/create/deployment").
				ContentType("application/json").
				Body(bytes.NewReader(body)).
				GetString(&response).
				Execute()
			if err != nil {
				fmt.Println("[Fault] " + err.Error())
			}
			fmt.Printf("create_deployment is %s and response is: %s\n", status, response)
		case yaml.Autoscaler_t:
			// 用来创建AutoScaler，需要发送给apiserver的参数为 autoScaler (def.AutoScaler)
			autoscaler, _ := yaml.ReadAutoScalerConfig(dir)
			request := *autoscaler
			response := ""
			body, _ := json.Marshal(request)
			err, status := httpget.Post("http://" + util.GetLocalIP().String() + ":" + fmt.Sprintf("%d", def.MasterPort) + "/create/autoscaler").
				ContentType("application/json").
				Body(bytes.NewReader(body)).
				GetString(&response).
				Execute()
			if err != nil {
				fmt.Println("[Fault] " + err.Error())
			}
			fmt.Printf("create_autoscaler is %s and response is: %s\n", status, response)
		case yaml.Gpujob_t:
			// 用来创建GPUJob，需要发送给apiserver的参数为 gpu (def.GPUJob)
			gpu, _ := yaml.ReadGPUJobConfig(dir)
			request := *gpu
			response := ""
			body, _ := json.Marshal(request)
			err, status := httpget.Post("http://" + util.GetLocalIP().String() + ":" + fmt.Sprintf("%d", def.MasterPort) + "/create/gpuJob").
				ContentType("application/json").
				Body(bytes.NewReader(body)).
				GetString(&response).
				Execute()
			if err != nil {
				fmt.Println("[Fault] " + err.Error())
			}
			fmt.Printf("create_gpuJob is %s and response is: %s\n", status, response)
		case yaml.Activity_t:
			// 用来创建function，需要发送给apiserver的参数为 function (def.Function)
			// 需要注意的是，返回的 response 会提供一个url，kubectl将它呈现给用户，后续用户可以使用它发送请求
			function, _ := yaml.ReadFunctionConfig(dir)
			request := *function
			response := ""
			body, _ := json.Marshal(request)
			err, status := httpget.Post("http://" + util.GetLocalIP().String() + ":" + fmt.Sprintf("%d", def.MasterPort) + "/create/function").
				ContentType("application/json").
				Body(bytes.NewReader(body)).
				GetString(&response).
				Execute()
			if err != nil {
				fmt.Println("[Fault] " + err.Error())
			}
			fmt.Printf("create_function is %s and response is: %s\n", status, response)
		}
	}

}
