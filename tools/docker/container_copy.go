package docker

import (
	"archive/tar"
	"bytes"
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

// CopyToContainer containerID:容器ID，
//dstPath: 目标文件在容器中的路径（文件夹，如/home/zuul/src/main/resources）,
//fileName:目标文件名（如application.yaml）,
//fileContent: 要写入文件的内容
func CopyToContainer(containerID string, dstPath string, fileName string, fileContent string) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		fmt.Printf("%v\n", err)
		panic(err)
	}
	defer cli.Close()

	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)
	err = tw.WriteHeader(&tar.Header{
		Name: fileName,
		Mode: 0777,
		Size: int64(len(fileContent)),
	})
	if err != nil {
		fmt.Printf("docker copy: %v\n", err)
	}
	tw.Write([]byte(fileContent))
	tw.Close()

	err = cli.CopyToContainer(context.Background(), containerID, dstPath, &buf, types.CopyToContainerOptions{
		AllowOverwriteDirWithFile: true,
	})
	if err != nil {
		fmt.Printf("%v\n", err)
		panic(err)
	} else {
		fmt.Printf("container %s has been removed\n", containerID)
	}
}
