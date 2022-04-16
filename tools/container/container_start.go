package container

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func StartContainer(containerID string) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		fmt.Printf("%v\n", err)
		panic(err)
	}
	defer cli.Close()
	err = cli.ContainerStart(context.Background(), containerID, types.ContainerStartOptions{})
	if err != nil {
		fmt.Printf("%v\n", err)
		panic(err)
	} else {
		fmt.Println("container", containerID, "start successfully")
	}
}
