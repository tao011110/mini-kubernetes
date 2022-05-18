package docker

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

// CommitContainer 需要传入参数为容器的ID
func CommitContainer(containerID string) string {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		fmt.Printf("%v\n", err)
		panic(err)
	}
	defer cli.Close()
	containerDetail, err := InspectContainer(containerID)
	if err != nil {
		fmt.Printf("%v\n", err)
		panic(err)
	}

	config := containerDetail.Config
	containerCommitOptions := types.ContainerCommitOptions{
		Config: config,
	}
	resp, err := cli.ContainerCommit(context.Background(), containerID, containerCommitOptions)
	if err != nil {
		fmt.Printf("%v\n", err)
		panic(err)
	} else {
		fmt.Println("Container ", containerID, "commit successfully")
	}
	fmt.Printf("resp.ID is %s\n", resp.ID)
	fmt.Printf("Image ID is %s\n", resp.ID[7:])

	return resp.ID[7:]
}
