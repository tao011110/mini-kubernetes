package create_api

import (
	"encoding/json"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"mini-kubernetes/tools/def"
	"mini-kubernetes/tools/etcd"
)

func CreatePod(cli *clientv3.Client, pod def.Pod) int {
	podInstance := def.PodInstance{}
	podInstance.Pod = pod

	//TODO: This node should be decided by scheduler in the future
	nodeID := 1
	nodeName := "node1"

	// 将新创建的pod写入到etcd当中
	podKey := "/pod/" + pod.Metadata.Name
	podValue, err := json.Marshal(pod)
	if err != nil {
		fmt.Printf("%v\n", err)
		panic(err)
	}
	etcd.Put(cli, podKey, string(podValue))

	//将新创建的podInstance写入到etcd当中
	podInstance.NodeID = uint64(nodeID)
	podInstanceKey := "/podInstance/" + pod.Metadata.Name
	podInstanceValue, err := json.Marshal(podInstance)
	if err != nil {
		fmt.Printf("%v\n", err)
		panic(err)
	}
	etcd.Put(cli, podInstanceKey, string(podInstanceValue))

	//更新相应node中的PodInstances列表
	nodeKey := "/node/" + nodeName
	nodeValue := etcd.Get(cli, nodeKey).Kvs[0].Value
	var node def.Node
	err = json.Unmarshal(nodeValue, &node)
	if err != nil {
		fmt.Printf("%v\n", err)
		panic(err)
	}
	node.PodInstances = append(node.PodInstances, podInstance)
	nodeValue, err = json.Marshal(node)
	if err != nil {
		fmt.Printf("%v\n", err)
		panic(err)
	}
	etcd.Put(cli, nodeKey, string(nodeValue))

	//更新kubelet watch的node-PodInstance table
	nodePIKey := "/nodePodInstance/" + nodeName
	nodePIValue := make([]byte, 0)
	podInstanceList := make([]string, 0)
	kvs := etcd.Get(cli, nodePIKey).Kvs
	if len(kvs) != 0 {
		nodePIValue = etcd.Get(cli, nodePIKey).Kvs[0].Value
		err = json.Unmarshal(nodePIValue, &podInstanceList)
		if err != nil {
			fmt.Printf("%v\n", err)
			panic(err)
		}
	}
	podInstanceList = append(podInstanceList, podInstanceKey)
	nodePIValue, err = json.Marshal(podInstanceList)
	if err != nil {
		fmt.Printf("%v\n", err)
		panic(err)
	}
	etcd.Put(cli, nodePIKey, string(nodePIValue))

	return nodeID
}
