package docker

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"mini-kubernetes/tools/def"
)

func CreateContainer(con def.Container, containerName string) string {
	ImageEnsure(con.Image)

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		fmt.Printf("%v\n", err)
		panic(err)
	}
	defer func(cli *client.Client) {
		_ = cli.Close()
	}(cli)

	config := GenerateConfig(con)
	body, err := cli.ContainerCreate(context.Background(), config, &container.HostConfig{},
		&network.NetworkingConfig{}, nil, containerName)
	if err != nil {
		fmt.Printf("%v\n", err)
		panic(err)
	}

	return body.ID
}
