package docker

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func ListNetwork() []types.NetworkResource {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		fmt.Printf("%v\n", err)
		panic(err)
	}
	defer cli.Close()
	networks, err := cli.NetworkList(context.Background(), types.NetworkListOptions{})
	if err != nil {
		fmt.Printf("%v\n", err)
		panic(err)
	}

	return networks
}
