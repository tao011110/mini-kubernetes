package kubelet_routines

import (
	"encoding/json"
	"mini-kubernetes/tools/def"
	"mini-kubernetes/tools/etcd"
	"mini-kubernetes/tools/pod"
	"mini-kubernetes/tools/util"
	"sort"
)

func EtcdWatcher(node *def.Node) {
	key := def.PodInstanceListKeyOfNode(node)
	watch := etcd.Watch(node.EtcdClient, key)
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
	sort.Strings(instancesNew)
	for index, instance := range node.PodInstances {
		if sort.SearchStrings(instancesNew, instance.ID) == len(instancesNew) {
			deleted = append(deleted, index)
		}
		instancesCurrent = append(instancesCurrent, instance.ID)
	}
	sort.Strings(instancesCurrent)
	for _, instance := range instancesNew {
		if sort.SearchStrings(instancesCurrent, instance) == len(instancesCurrent) {
			added = append(added, instance)
		}
	}
	return
}
