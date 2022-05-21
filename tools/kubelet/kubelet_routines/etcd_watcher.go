package kubelet_routines

import (
	"encoding/json"
	"fmt"
	"mini-kubernetes/tools/def"
	"mini-kubernetes/tools/etcd"
	"mini-kubernetes/tools/pod"
	"mini-kubernetes/tools/util"
)

func EtcdWatcher(node *def.Node) {
	key := def.PodInstanceListKeyOfNode(node)
	watch := etcd.Watch(node.EtcdClient, key)
	fmt.Println("watch now!")
	for wc := range watch {
		for _, w := range wc.Events {
			var instances []string
			_ = json.Unmarshal(w.Kv.Value, &instances)
			handlePodInstancesChange(node, instances)
		}
	}
}

func handlePodInstancesChange(node *def.Node, instances []string) {
	adds, deletes := comparePodList(node, instances)
	fmt.Println("add are", adds)
	fmt.Println("deletes are", deletes)
	for _, add := range adds {
		podInstance := util.GetPodInstance(add, node.EtcdClient)
		if podInstance.Status != def.PENDING {
			continue
		}
		pod.CreateAndStartPod(&podInstance, node)
		node.PodInstances = append(node.PodInstances, &podInstance)
	}
	for _, index := range deletes {
		if node.PodInstances[index].Status != def.RUNNING {
			continue
		}
		pod.StopAndRemovePod(node.PodInstances[index], node)
		//node.PodInstances = append(node.PodInstances[:index], node.PodInstances[index+1:]...)
	}
}

func comparePodList(node *def.Node, instancesNew []string) (added []string, deleted []int) {
	var instancesCurrent []string
	for _, instance := range node.PodInstances {
		instancesCurrent = append(instancesCurrent, instance.ID)
	}
	fmt.Println("instancesCurrent  ", instancesCurrent)
	fmt.Println("instancesNew  ", instancesNew)
	added, deletedIDs := util.DifferTwoStringList(instancesCurrent, instancesNew)
	for _, delete_ := range deletedIDs {
		for index, instance := range instancesCurrent {
			if delete_ == instance {
				deleted = append(deleted, index)
				break
			}
		}
	}
	return
}
