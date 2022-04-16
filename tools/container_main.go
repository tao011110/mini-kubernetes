package main

import (
	"fmt"
	"github.com/YOUR-USER-OR-ORG-NAME/YOUR-REPO-NAME/tools/container"
)

func main() {
	path := "./container/container_test2.yaml"

	containerIDs := container.CreateContrainer(path)

	for _, id := range containerIDs {
		fmt.Printf("has created: %s\n", id)
		container.StartContainer(id)
		//container.RestartContainer(id)
		//container.StopContainer(id)
		//container.RemoveContainer(id)
	}
}
