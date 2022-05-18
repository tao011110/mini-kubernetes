package docker_test

import (
	//"mini-kubernetes/tools/docker"
	"mini-kubernetes/tools/docker"
	"testing"
)

func Test(t *testing.T) {
	//path := "./docker_test3.yaml"

	//containerIDs := docker.CreateContrainer(path, "10.24.0.0")

	//for _, id := range containerIDs {
	//	t.Logf("has created: %s\n", id)
	//	docker.StartContainer(id)
	//	//container.RestartContainer(id)
	//	//container.StopContainer(id)
	//	//container.RemoveContainer(id)
	//}
	imageID := docker.CommitContainer("2dbd1987a731")
	docker.TagImage(imageID, "test1")
	docker.PushImage("test1")

	t.Log("test finished\n")
}
