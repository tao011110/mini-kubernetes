package container

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"io"
	"io/ioutil"
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
		resp, err := cli.ImagePull(context.Background(), targetImage, types.ImagePullOptions{})
		if err != nil {
			fmt.Printf("Pull image failed %s\n %v\n", targetImage, err)
			panic(err)
		}
		_, err = io.Copy(ioutil.Discard, resp)
		if err != nil {
			fmt.Printf("%v\n", err)
			panic(err)
		} else {
			fmt.Printf("Pulled image %s successully\n", targetImage)
		}
	}
}
