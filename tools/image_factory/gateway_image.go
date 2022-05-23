package image_factory

import (
	"fmt"
	"mini-kubernetes/tools/def"
	"mini-kubernetes/tools/docker"
	"mini-kubernetes/tools/gateway"
	"os/exec"
)

func MakeGatewayImage(dns *def.DNSDetail, nameGatewayImageName string) {
	container := def.Container{
		Image: def.PyFunctionTemplateImage,
	}
	containerID := docker.CreateContainer(container, nameGatewayImageName)
	fileStr := gateway.GenerateApplicationYaml(*dns)
	docker.CopyToContainer(containerID, def.PyHandlerParentDirPath, def.PyHandlerFileName, fileStr)
	cmd := exec.Command("docker", "exec", containerID, "/bin/bash", "-c", fmt.Sprintf("'%s'", def.GatewayPackageCmd)).String()
	WriteCmdToFile(def.TemplateCmdFilePath, cmd)
	command := fmt.Sprintf(`%s .`, def.TemplateCmdFilePath)
	err := exec.Command("/bin/bash", "-c", command).Run()
	if err != nil {
		fmt.Println(err)
	}
	docker.CommitContainer(containerID, nameGatewayImageName)
	docker.PushImage(nameGatewayImageName)
	docker.StopContainer(containerID)
	_, _ = docker.RemoveContainer(containerID)
}
