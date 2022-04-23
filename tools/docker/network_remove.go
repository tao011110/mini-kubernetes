package docker

import (
	"context"
	"fmt"
	"github.com/docker/docker/client"
)

func RemoveNetBridge(networkID string) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		fmt.Printf("%v\n", err)
		panic(err)
	}
	defer cli.Close()
	err = cli.NetworkRemove(context.Background(), networkID)
	if err != nil {
		fmt.Printf("%v\n", err)
		panic(err)
	} else {
		fmt.Printf("bridge %s has been removed\n", networkID)
	}
}
