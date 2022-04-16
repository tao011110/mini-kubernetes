package main

import (
	"fmt"
	"github.com/YOUR-USER-OR-ORG-NAME/YOUR-REPO-NAME/tools/docker"
)

func main() {
	path := "./container/docker_test3.yaml"

	containerIDs := docker.CreateContrainer(path)

	for _, id := range containerIDs {
		fmt.Printf("has created: %s\n", id)
		docker.StartContainer(id)
		//container.RestartContainer(id)
		//container.StopContainer(id)
		//container.RemoveContainer(id)
	}
}
