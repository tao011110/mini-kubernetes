package docker_test

import (
	"mini-kubernetes/tools/def"
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

	//用户生成image并 push到仓库中，需要提供函数名 funcName 和用于打包成镜像的容器ID containerID
	//funcName := "test1"
	//containerID := "c2d46cf5ce63"
	//docker.CommitContainer(containerID, funcName)
	//docker.PushImage(funcName)

	container := def.Container{
		Image: "registry.cn-hangzhou.aliyuncs.com/taoyucheng/mink8s:tmpForGateway",
	}
	docker.CreateContainer(container, "test")

	t.Log("test finished\n")
}
