package pod

import (
	"mini-kubernetes/tools/docker"
)

func (podInstance *PodInstance) RestartPod() {
	podInstance.Status = RESTARTING
	for index, container := range podInstance.ContainerStatus {
		docker.RestartContainer(container.ID)
		podInstance.ContainerStatus[index].Status = SUCCEEDED
	}
	podInstance.Status = SUCCEEDED
	podInstance.PodInstanceStatus = InstanceStatus{}
}
