package docker_test

import (
	"mini-kubernetes/tools/docker"
	"testing"
)

func Test(t *testing.T) {
	path := "./docker_test3.yaml"

	containerIDs := docker.CreateContrainer(path, "172.18.0.0")

	for _, id := range containerIDs {
		t.Logf("has created: %s\n", id)
		docker.StartContainer(id)
		//container.RestartContainer(id)
		//container.StopContainer(id)
		//container.RemoveContainer(id)
	}
	t.Log("test finished\n")
}
