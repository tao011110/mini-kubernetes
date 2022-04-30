package pod

import (
	"mini-kubernetes/tools/def"
	"mini-kubernetes/tools/docker"
	"mini-kubernetes/tools/util"
)

func (podInstance *PodInstance) StopAndRemovePod(node *def.Node) {
	podInstance.Status = SUCCEEDED
	util.PersistPodInstance(*podInstance, node.EtcdClient)
	for index, container := range podInstance.ContainerSpec {
		docker.StopContainer(container.ID)
		_, _ = docker.RemoveContainer(container.ID)
		podInstance.ContainerSpec[index].Status = SUCCEEDED
		util.PersistPodInstance(*podInstance, node.EtcdClient)
	}
}
