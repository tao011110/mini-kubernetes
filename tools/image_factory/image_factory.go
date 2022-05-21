package image_factory

import (
	"fmt"
	"mini-kubernetes/tools/def"
	"mini-kubernetes/tools/docker"
	"os/exec"
	"strings"
)

func EchoFactory(preStr string, target string) string {
	newStr := strings.Replace(preStr, "\n", "\\n", -1)
	newStr = strings.Replace(newStr, "\\", "\\\\", -1)
	return fmt.Sprintf("echo -e \\\"%s\\\"", newStr)
}

func ImageFactory(baseImageName string, newImageName string, commandInContainer []string) {
	container := def.Container{
		Image: baseImageName,
	}
	containerID := docker.CreateContainer(container, newImageName)
	for _, command := range commandInContainer {
		cmd := exec.Command("docker", "exec", containerID, "/bin/bash", "-c", fmt.Sprintf("'%s'", command))
		fmt.Println(cmd.String())
		err := cmd.Run()
		if err != nil {
			fmt.Println(err)
		}
	}
	docker.CommitContainer(containerID, newImageName)
	docker.PushImage(newImageName)
	docker.StopContainer(containerID)
	_, _ = docker.RemoveContainer(containerID)
}
