package image_factory

import (
	"mini-kubernetes/tools/def"
	"mini-kubernetes/tools/docker"
	"mini-kubernetes/tools/gateway"
)

func MakeGatewayImage(dns *def.DNSDetail, nameGatewayImageName string) {
	container := def.Container{
		Image: def.GatewayImage,
	}
	containerID := docker.CreateContainer(container, nameGatewayImageName)
	fileStr := gateway.GenerateApplicationYaml(*dns)
	docker.CopyToContainer(containerID, def.RequirementsParentDirPath, def.GatewayRoutesConfigPathInImage, fileStr)
	//cmd := exec.Command("docker", "exec", containerID, "/bin/bash", "-c", fmt.Sprintf("'%s'", def.GatewayPackageCmd)).String()
	//WriteCmdToFile(def.TemplateCmdFilePath, cmd)
	//command := fmt.Sprintf(`%s .`, def.TemplateCmdFilePath)
	//err := exec.Command("/bin/bash", "-c", command).Run()
	//if err != nil {
	//	fmt.Println(err)
	//}
	docker.CommitContainer(containerID, nameGatewayImageName)
	docker.PushImage(nameGatewayImageName)
	docker.StopContainer(containerID)
	_, _ = docker.RemoveContainer(containerID)
}
