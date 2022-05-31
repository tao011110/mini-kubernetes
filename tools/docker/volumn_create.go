package docker

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
)

func CreateVolume(containerID string) (types.Volume, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		fmt.Printf("%v\n", err)
		panic(err)
	}
	defer cli.Close()
	options := volume.VolumeCreateBody{}
	volume, err := cli.VolumeCreate(context.Background(), options)
	if err != nil {
		fmt.Printf("%v\n", err)
		panic(err)
	}

	return volume, err
}
