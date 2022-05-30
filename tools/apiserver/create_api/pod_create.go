package create_api

import (
	"encoding/json"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"mini-kubernetes/tools/def"
	"mini-kubernetes/tools/etcd"
	"strconv"
)

func CreatePod(cli *clientv3.Client, pod_ def.Pod) def.PodInstance {
	podInstance := def.PodInstance{}
	podInstance.Pod = pod_

	//TODO:This node should be decided by scheduler in the future
	//nodeID := 1
	//node_test := def.Node{
	//	NodeID: nodeID,
	//}

	// 将新创建的pod写入到etcd当中
	podKey := "/pod/" + pod_.Metadata.Name
	podValue, err := json.Marshal(pod_)
	if err != nil {
		fmt.Printf("%v\n", err)
		panic(err)
	}
	etcd.Put(cli, podKey, string(podValue))

	//将新创建的podInstance写入到etcd当中
	//podInstance.NodeID = nodeID
	podInstanceKey := def.GetKeyOfPodInstance(pod_.Metadata.Name)
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

	//更新PodInstanceIDList
	podInstanceIDList := make([]string, 0)
	kvs := etcd.Get(cli, def.PodInstanceListID).Kvs
	if len(kvs) != 0 {
		podInstanceIDListValue := kvs[0].Value
		err := json.Unmarshal(podInstanceIDListValue, &podInstanceIDList)
		if err != nil {
			fmt.Printf("%v\n", err)
			panic(err)
		}
	}
	podInstanceIDList = append(podInstanceIDList, podInstance.ID)
	podInstanceIDValue, err := json.Marshal(podInstanceIDList)
	if err != nil {
		fmt.Printf("%v\n", err)
		panic(err)
	}
	etcd.Put(cli, def.PodInstanceListID, string(podInstanceIDValue))

	//更新相应node中的PodInstances列表
	//nodeKey := "/node/" + strconv.Itoa(nodeID)
	//nodeValue := etcd.Get(cli, nodeKey).Kvs[0].Value
	//var node def.Node
	//err = json.Unmarshal(nodeValue, &node)
	//if err != nil {
	//	fmt.Printf("%v\n", err)
	//	panic(err)
	//}
	//node.PodInstances = append(node.PodInstances, &podInstance)
	//nodeValue, err = json.Marshal(node)
	//if err != nil {
	//	fmt.Printf("%v\n", err)
	//	panic(err)
	//}
	//etcd.Put(cli, nodeKey, string(nodeValue))

	//更新kubelet watch的node-PodInstance table
	//nodePIKey := def.PodInstanceListKeyOfNode(&node_test)
	//nodePIValue := make([]byte, 0)
	//podInstanceList := make([]string, 0)
	//kvs := etcd.Get(cli, nodePIKey).Kvs
	//if len(kvs) != 0 {
	//	nodePIValue = etcd.Get(cli, nodePIKey).Kvs[0].Value
	//	err = json.Unmarshal(nodePIValue, &podInstanceList)
	//	if err != nil {
	//		fmt.Printf("%v\n", err)
	//		panic(err)
	//	}
	//}
	//podInstanceList = append(podInstanceList, podInstanceKey)
	//nodePIValue, err = json.Marshal(podInstanceList)
	//if err != nil {
	//	fmt.Printf("%v\n", err)
	//	panic(err)
	//}
	//etcd.Put(cli, nodePIKey, string(nodePIValue))

	return podInstance
}

// CheckAddInService 检查是否需要将新创建的pod纳入到之前创建好的service当中
func CheckAddInService(cli *clientv3.Client, podInstance def.PodInstance) []def.Service {
	servicePrefix := "/service/"
	kvs := etcd.GetWithPrefix(cli, servicePrefix).Kvs
	serviceList := make([]def.Service, 0)
	for _, kv := range kvs {
		service := def.Service{}
		serviceValue := make([]byte, 0)
		serviceValue = kv.Value
		err := json.Unmarshal(serviceValue, &service)
		if err != nil {
			fmt.Printf("%v\n", err)
			panic(err)
		}
		// 匹配发现该pod应该被纳入到service的监管下
		fmt.Printf("in pod create podInstance.Pod.Metadata.Label is %s\n", podInstance.Pod.Metadata.Label)
		fmt.Printf("in pod create service.Selector.Name is %s\n", service.Selector.Name)
		//tmpBindings := make([]def.PortsBindings, 0)
		if podInstance.Pod.Metadata.Label == service.Selector.Name {
			//for _, binding := range service.PortsBindings {
			//	tmpBinding := def.PortsBindings{
			//		Ports: binding.Ports,
			//	}
			//	tmpEndpoints := make([]string, 0)
			//	for _, container := range podInstance.Pod.Spec.Containers {
			//		for _, portMapping := range container.PortMappings {
			//			containerPort := strconv.Itoa(int(portMapping.ContainerPort))
			//			fmt.Println(binding.Ports.TargetPort)
			//			fmt.Println(containerPort)
			//			if binding.Ports.TargetPort == containerPort {
			//				tmpEndpoints = append(tmpEndpoints, podInstance.IP+":"+binding.Ports.TargetPort)
			//				fmt.Printf("podInstance.IP is %v\n", podInstance.IP)
			//			}
			//		}
			//	}
			//	tmpBinding.Endpoints = tmpEndpoints
			//	tmpBindings = append(tmpBindings, tmpBinding)
			//}
			service.PortsBindings = AddPodInstanceIntoService(podInstance, service)
			fmt.Printf("now service si %v\n", service)
			serviceList = append(serviceList, service)
		}
	}

	// 若pod被加入到service当中，将service重新写入到etcd中
	for _, svc := range serviceList {
		svcKey := "/service/" + svc.Name
		svcValue, err := json.Marshal(svc)
		if err != nil {
			fmt.Printf("%v\n", err)
			panic(err)
		}
		etcd.Put(cli, svcKey, string(svcValue))
	}

	return serviceList
}

func AddPodInstanceIntoService(podInstance def.PodInstance, service def.Service) []def.PortsBindings {
	tmpBindings := make([]def.PortsBindings, 0)
	for _, binding := range service.PortsBindings {
		tmpBinding := def.PortsBindings{
			Ports: binding.Ports,
		}
		tmpEndpoints := binding.Endpoints
		for _, container := range podInstance.Pod.Spec.Containers {
			for _, portMapping := range container.PortMappings {
				containerPort := strconv.Itoa(int(portMapping.ContainerPort))
				fmt.Println(binding.Ports.TargetPort)
				fmt.Println(containerPort)
				if binding.Ports.TargetPort == containerPort {
					tmpEndpoints = append(tmpEndpoints, podInstance.IP+":"+binding.Ports.TargetPort)
					fmt.Printf("podInstance.IP is %v\n", podInstance.IP)
				}
			}
		}
		tmpBinding.Endpoints = tmpEndpoints
		tmpBindings = append(tmpBindings, tmpBinding)
	}

	return tmpBindings
}
