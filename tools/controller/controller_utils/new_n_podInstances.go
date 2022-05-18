package controller_utils

import (
	clientv3 "go.etcd.io/etcd/client/v3"
	"mini-kubernetes/tools/def"
)

func NewNPodInstance(etcdClient *clientv3.Client, podName string, num int) {
	pod := GetPodByName(etcdClient, podName)
	for i := 0; i < num; i++ {
		podInstance := def.PodInstance{
			Pod:           *pod,
			ID:            def.GenerateKeyOfPodInstanceReplicas(podName),
			NodeID:        def.NodeUndefined,
			Status:        def.PENDING,
			ContainerSpec: make([]def.ContainerStatus, len(pod.Spec.Containers)),
			RestartCount:  0,
		}
		AddPodInstance(etcdClient, &podInstance)
	}
}
