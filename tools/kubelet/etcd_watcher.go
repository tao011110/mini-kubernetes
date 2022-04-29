package kubelet

import (
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	clientv3 "go.etcd.io/etcd/client/v3"
	"mini-kubernetes/tools/def"
	"mini-kubernetes/tools/etcd"
	"mini-kubernetes/tools/pod"
	"sort"
)

type PodInstances struct {
	Instances []string
}

func EtcdWatcher(cli *clientv3.Client, nodeID int, e *echo.Echo) {
	key := fmt.Sprintf("Node%d_podInstances", nodeID)
	watch := etcd.Watch(cli, key)
	for wc := range watch {
		for _, w := range wc.Events {
			var instances []string
			_ = json.Unmarshal(w.Kv.Value, &instances)
			handlePodInstancesChange(instances, e, cli)
		}
	}
}

func handlePodInstancesChange(instances []string, e *echo.Echo, cli *clientv3.Client) {
	adds, deletes := comparePodList(instances)
	for _, add := range adds {
		resp := etcd.Get(cli, add)
		podInstance := def.PodInstance{}
		jsonString := ``
		for _, ev := range resp.Kvs {
			jsonString += fmt.Sprintf(`"%s":"%s"`, ev.Key, ev.Value)
		}
		jsonString = fmt.Sprintf(`{%s}`, jsonString)
		_ = json.Unmarshal([]byte(jsonString), &podInstance)
		pod.CreateAndStartPod(&podInstance)
		node.PodInstances = append(node.PodInstances, podInstance)
	}
	for _, index := range deletes {
		pod.StopAndRemovePod(&node.PodInstances[index])
	}
}

func comparePodList(instancesNew []string) (added []string, deleted []int) {
	var instancesCurrent []string
	sort.Strings(instancesNew)
	for index, instance := range node.PodInstances {
		if sort.SearchStrings(instancesNew, instance.Name) == len(instancesNew) {
			deleted = append(deleted, index)
		}
		instancesCurrent = append(instancesCurrent, instance.Name)
	}
	sort.Strings(instancesCurrent)
	for _, instance := range instancesNew {
		if sort.SearchStrings(instancesCurrent, instance) == len(instancesCurrent) {
			added = append(added, instance)
		}
	}
	return
}
