package docker_test

import (
	"fmt"
	"github.com/YOUR-USER-OR-ORG-NAME/YOUR-REPO-NAME/tools/docker"
	"testing"
)

func Test(t *testing.T) {
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
