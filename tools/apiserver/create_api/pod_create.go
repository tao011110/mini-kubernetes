package create_api

import (
	"encoding/json"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"mini-kubernetes/tools/def"
	"mini-kubernetes/tools/etcd"
	"strconv"
)

func CreatePod(cli *clientv3.Client, pod_ def.Pod) int {
	podInstance := def.PodInstance{}
	podInstance.Pod = pod_

	//TODO:This node should be decided by scheduler in the future
	nodeID := 1

	// 将新创建的pod写入到etcd当中
	podKey := "/pod/" + pod_.Metadata.Name
	podValue, err := json.Marshal(pod_)
	if err != nil {
		fmt.Printf("%v\n", err)
		panic(err)
	}
	etcd.Put(cli, podKey, string(podValue))

	//将新创建的podInstance写入到etcd当中
	podInstance.NodeID = nodeID
	podInstanceKey := "/podInstance/" + pod_.Metadata.Name
	podInstance.ID = podInstanceKey
	podInstance.ContainerSpec = make([]def.ContainerStatus, len(pod_.Spec.Containers))

	// NOTICE: 此处的一些设置仅供测试使用
	//podInstance.ClusterIP = "10.24.1.2"
	//podInstance.StartTime = time.Now()
	//podInstance.RestartCount = 0
	//podInstance.Status = def.RUNNING
	//for _, container := range podInstance.ContainerSpec {
	//	container.Status = def.RUNNING
	//}

	podInstanceValue, err := json.Marshal(podInstance)
	if err != nil {
		fmt.Printf("%v\n", err)
		panic(err)
	}
	etcd.Put(cli, podInstanceKey, string(podInstanceValue))

	//更新相应node中的PodInstances列表
	nodeKey := "/node/" + strconv.Itoa(nodeID)
	nodeValue := etcd.Get(cli, nodeKey).Kvs[0].Value
	var node def.Node
	err = json.Unmarshal(nodeValue, &node)
	if err != nil {
		fmt.Printf("%v\n", err)
		panic(err)
	}
	node.PodInstances = append(node.PodInstances, &podInstance)
	nodeValue, err = json.Marshal(node)
	if err != nil {
		fmt.Printf("%v\n", err)
		panic(err)
	}
	etcd.Put(cli, nodeKey, string(nodeValue))

	//更新kubelet watch的node-PodInstance table
	nodePIKey := "/nodePodInstance/" + strconv.Itoa(nodeID)
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
