package docker

import (
	"context"
	"fmt"
	"github.com/docker/docker/client"
	"mini-kubernetes/tools/def"
)

// TagImage 调用时需要传入容器的ID和函数名称
func TagImage(containerID string, funcName string) {
	cli, _ := client.NewClientWithOpts(client.FromEnv)
	err := cli.ImageTag(context.Background(), containerID, def.RgistryAddr+funcName)
	if err != nil {
		fmt.Printf("%v\n", err)
		panic(err)
	} else {
		fmt.Println("Container ", containerID, "commit successfully")
	}
}
