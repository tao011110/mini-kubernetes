package pod

import (
	"mini-kubernetes/tools/def"
	"mini-kubernetes/tools/docker"
)

func RestartPod(podInstance *def.PodInstance) {
	podInstance.Status = def.RESTARTING
	for index, container := range podInstance.ContainerSpec {
		docker.RestartContainer(container.ID)
		podInstance.ContainerSpec[index].Status = def.SUCCEEDED
	}
	podInstance.Status = def.SUCCEEDED
	//podInstance.PodInstanceStatus = def.InstanceSpec{}
}
