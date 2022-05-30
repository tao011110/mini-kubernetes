package docker

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func ImageEnsure(targetImage string) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		fmt.Printf("%v\n", err)
		panic(err)
	}
	defer cli.Close()

	images, err := cli.ImageList(context.Background(), types.ImageListOptions{})
	if err != nil {
		fmt.Printf("%v\n", err)
		panic(err)
	}
	isImageExist := false
	for _, image := range images {
		for _, tag := range image.RepoTags {
			if targetImage == tag {
				isImageExist = true
				break
			}
		}
	}
	if !isImageExist {
		fmt.Printf("Image %s doesn't exist locally, try to pull it now\n", targetImage)
	}
	PullImage(targetImage)
}
