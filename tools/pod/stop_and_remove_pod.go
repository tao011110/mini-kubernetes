package pod

import (
	"mini-kubernetes/tools/def"
	"mini-kubernetes/tools/docker"
)

func StopAndRemovePod(podInstance *def.PodInstance) {
	for index, container := range podInstance.ContainerStatus {
		docker.StopContainer(container.ID)
		_, _ = docker.RemoveContainer(container.ID)
		podInstance.ContainerStatus[index].Status = def.SUCCEEDED
	}
	podInstance.Status = def.SUCCEEDED
}
