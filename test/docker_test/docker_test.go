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

	//TODO: 用户生成image并 push到仓库中，需要提供函数名 funcName 和用于打包成镜像的容器ID containerID
	funcName := "test1"
	containerID := "2dbd1987a731"
	docker.CommitContainer(containerID, funcName)
	docker.PushImage(funcName)

	t.Log("test finished\n")
}
