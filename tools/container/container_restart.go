package container

import (
	"context"
	"fmt"
	"github.com/docker/docker/client"
	"time"
)

func RestartContainer(containerID string) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		fmt.Printf("%v\n", err)
		panic(err)
	}
	defer cli.Close()
	timeout := time.Second * 5
	err = cli.ContainerRestart(context.Background(), containerID, &timeout)
	if err != nil {
		fmt.Printf("%v\n", err)
		panic(err)
	} else {
		fmt.Printf("container %s has been restarted\n", containerID)
	}
}
