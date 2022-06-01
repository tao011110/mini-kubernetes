package pod

import (
	"mini-kubernetes/tools/def"
	"mini-kubernetes/tools/docker"
	"mini-kubernetes/tools/util"
)

func StopAndRemovePod(podInstance *def.PodInstance, node *def.Node) {
	if podInstance.Status == def.RUNNING {
		podInstance.Status = def.SUCCEEDED
	}
	util.PersistPodInstance(*podInstance, node.EtcdClient)
	for index, container := range podInstance.ContainerSpec {
		if podInstance.ContainerSpec[index].Status == def.RUNNING {
			podInstance.ContainerSpec[index].Status = def.SUCCEEDED
		}
		docker.StopContainer(container.ID)
		_, _ = docker.RemoveContainer(container.ID)
		util.PersistPodInstance(*podInstance, node.EtcdClient)
	}
	docker.StopContainer(podInstance.PauseContainer.ID)
	_, _ = docker.RemoveContainer(podInstance.PauseContainer.ID)

	if podInstance.PauseContainer.Status == def.RUNNING {
		podInstance.PauseContainer.Status = def.SUCCEEDED
	}
	util.PersistPodInstance(*podInstance, node.EtcdClient)
}
