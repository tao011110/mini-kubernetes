package pod

import (
	"mini-kubernetes/tools/docker"
)

func (podInstance *PodInstance) StopAndRemovePod() {
	podInstance.Status = SUCCEEDED
	for index, container := range podInstance.ContainerStatus {
		docker.StopContainer(container.ID)
		_, _ = docker.RemoveContainer(container.ID)
		podInstance.ContainerStatus[index].Status = SUCCEEDED
	}
}
