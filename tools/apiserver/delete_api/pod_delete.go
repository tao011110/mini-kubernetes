package delete_api

import (
	"encoding/json"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"mini-kubernetes/tools/def"
	"mini-kubernetes/tools/etcd"
	"strconv"
)

func DeletePod(cli *clientv3.Client, podInstanceName string) bool {
	//在etcd中删除podInstance
	podInstanceKey := "/podInstance/" + podInstanceName
	resp := etcd.Get(cli, podInstanceKey)
	if len(resp.Kvs) == 0 {
		return false
	}
	podInstanceValue := resp.Kvs[0].Value
	etcd.Delete(cli, podInstanceKey)

	//更新相应node中的PodInstances列表
	podInstance := def.PodInstance{}
	err := json.Unmarshal(podInstanceValue, &podInstance)
	if err != nil {
		fmt.Printf("%v\n", err)
		panic(err)
	}
	nodeKey := "/node/" + strconv.Itoa(int(podInstance.NodeID))
	nodeValue := etcd.Get(cli, nodeKey).Kvs[0].Value
	var node def.Node
	err = json.Unmarshal(nodeValue, &node)
	if err != nil {
		fmt.Printf("%v\n", err)
		panic(err)
	}
	podInstances := make([]*def.PodInstance, len(node.PodInstances)-1)
	podInstanceList := make([]string, len(node.PodInstances)-1)
	for _, pi := range node.PodInstances {
		if pi.Pod.Metadata.Name != podInstance.Pod.Metadata.Name {
			podInstances = append(podInstances, pi)
			podInstanceList = append(podInstanceList, "/nodePodInstance/"+pi.Metadata.Name)
		}
	}
	node.PodInstances = podInstances
	nodeValue, err = json.Marshal(node)
	if err != nil {
		fmt.Printf("%v\n", err)
		panic(err)
	}
	etcd.Put(cli, nodeKey, string(nodeValue))

	//更新kubelet watch的node-PodInstance table
	nodePIKey := "/nodePodInstance/" + strconv.Itoa(int(podInstance.NodeID))
	nodePIValue, err := json.Marshal(podInstanceList)
	if err != nil {
		fmt.Printf("%v\n", err)
		panic(err)
	}
	etcd.Put(cli, nodePIKey, string(nodePIValue))

	return true
}
