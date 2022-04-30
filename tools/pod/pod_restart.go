package pod

import (
	"mini-kubernetes/tools/docker"
)

func (podInstance *PodInstance) RestartPod() {
	podInstance.Status = RESTARTING
	for index, container := range podInstance.ContainerSpec {
		docker.RestartContainer(container.ID)
		podInstance.ContainerSpec[index].Status = SUCCEEDED
	}
	podInstance.Status = SUCCEEDED
	podInstance.PodInstanceStatus = InstanceSpec{}
}
