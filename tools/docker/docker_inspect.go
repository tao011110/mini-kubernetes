package docker

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func InspectContainer(containerID string) (types.ContainerJSON, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		fmt.Printf("%v\n", err)
		panic(err)
	}
	defer cli.Close()
	status, err := cli.ContainerInspect(context.Background(), containerID)
	if err != nil {
		fmt.Printf("%v\n", err)
		panic(err)
	}

	return status, err
}
