package docker_test

import (
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

	//container := def.Container{
	//	Image: "registry.cn-hangzhou.aliyuncs.com/taoyucheng/mink8s:tmpForGateway",
	//}
	//docker.CreateContainer(container, "test")

	content := "server:\n  port: 80\nspring:\n  application:\n    name: zuul\nzuul:\n  routes:\n    route0:\n      path: /route1/**\n      url: http://192.168.40.10:300\n    route1:\n      path: /route2/**\n      url: http://192.168.40.10:80\n"
	docker.CopyToContainer("0d29cc77b9f857d3dfb4363d4cb7ca5545bda7f9b5572397940cf4d7394ca5ff",
		"/home/zuul/src/main/resources", "application.yaml", content)

	t.Log("test finished\n")
}
