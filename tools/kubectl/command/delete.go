package command

import (
	"fmt"
	"mini-kubernetes/tools/def"
	"mini-kubernetes/tools/util"

	"mini-kubernetes/tools/httpget"
	"mini-kubernetes/tools/yaml"

	"github.com/urfave/cli"
)

func NewDeleteCommand() cli.Command {
	return cli.Command{
		Name:  "delete",
		Usage: "Delete resources according to xxx.yaml or name",
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
	/* kubectl delete 通过配置文件名、资源名称来删除资源。 */
	dir := c.String("file")
	var src_type int
	var src_name string
	if dir != "" {
		// 根据yaml文件来删除资源
		fmt.Printf("Using dir: %s\n", dir)
		src_type, src_name, _ = yaml.ReadTypeAndName(dir)
	} else {
		// 根据名称来删除资源
		if len(c.Args()) == 0 {
			fmt.Println("You need to specify pod or service")
			return
		}
		if c.Args()[0] == "pod" {
			src_type = yaml.Pod_t
			src_name = c.Args()[1]
			fmt.Printf("Delete pod whose name is : %s\n", src_name)
		} else if c.Args()[0] == "service" {
			// 目前还不能确定是哪种service类型
			src_type = yaml.Unknown_t
			src_name = c.Args()[1]
			fmt.Printf("Delete service whose name is : %s\n", src_name)
		} else if c.Args()[0] == "deployment" {
			src_type = yaml.Deployment_t
			src_name = c.Args()[1]
			fmt.Printf("Delete deployment whose name is : %s\n", src_name)
		}
	}

	if src_type >= yaml.Pod_t && src_type <= yaml.Unknown_t && src_name != "" {
		if src_type == yaml.Pod_t {
			// 格式 kubectl delete pod xxx
			// 需要发送给apiserver的参数为 podName string
			response := ""
			err, status := httpget.DELETE("http://" + util.GetLocalIP().String() + ":" + fmt.Sprintf("%d", def.MasterPort) + "/delete_pod/" + src_name).
				ContentType("application/json").
				GetString(&response).
				Execute()
			if err != nil {
				fmt.Println("[Fault] " + err.Error())
			} else {
				fmt.Printf("get_pod status is %s\n", status)
				if status == "200" {
					fmt.Printf("Delete pod %s successfully and the response is: %v\n", src_name, response)
				} else {
					fmt.Printf("Pod %s doesn't exist\n", src_name)
				}
			}
		} else if src_type == yaml.ClusterIP_t || src_type == yaml.Nodeport_t || src_type == yaml.Unknown_t {
			// 格式 kubectl delete service xxx
			// 需要发送给apiserver的参数为 serviceName string
			response := ""
			err, status := httpget.DELETE("http://" + util.GetLocalIP().String() + ":" + fmt.Sprintf("%d", def.MasterPort) + "/delete/service/" + src_name).
				ContentType("application/json").
				GetString(&response).
				Execute()
			if err != nil {
				fmt.Println("[Fault] " + err.Error())
			}
			fmt.Printf("delete clusterIPService status is %s\n", status)
			if status == "200" {
				fmt.Printf("Delete service %s successfully and the response is: %v\n", src_name, response)
			} else {
				fmt.Printf("Service %s doesn't exist\n", src_name)
			}
		} else if src_type == yaml.Deployment_t {
			// 格式 kubectl delete deployment xxx
			// 用来删除deployment，需要发送给apiserver的参数为 deploymentName(string)
			response := ""
			err, status := httpget.DELETE("http://" + util.GetLocalIP().String() + ":" + fmt.Sprintf("%d", def.MasterPort) + "/delete/deployment/" + src_name).
				ContentType("application/json").
				GetString(&response).
				Execute()
			if err != nil {
				fmt.Println("[Fault] " + err.Error())
			}
			fmt.Printf("delete deployment status is %s\n", status)
			if status == "200" {
				fmt.Printf("delete deployment %s successfully and the response is: %v\n", src_name, response)
			} else {
				fmt.Printf("deployment %s doesn't exist\n", src_name)
			}
		} else {
			fmt.Println("Now delete only support pod/service/deployment")
		}
	}

}
