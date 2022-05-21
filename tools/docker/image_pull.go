package docker

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"io"
	"io/ioutil"
)

func PullImage(targetImage string) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		fmt.Printf("%v\n", err)
		panic(err)
	}
	defer cli.Close()

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
