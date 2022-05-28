package delete_api

import (
	"encoding/json"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"mini-kubernetes/tools/def"
	"mini-kubernetes/tools/etcd"
	"strings"
)

func DeletePod(cli *clientv3.Client, podInstanceName string) (bool, def.PodInstance) {
	//在etcd中删除podInstance
	podInstanceKey := "/podInstance/" + podInstanceName
	podInstance := def.PodInstance{}
	resp := etcd.Get(cli, podInstanceKey)
	if len(resp.Kvs) == 0 {
		return false, podInstance
	}
	podInstanceValue := resp.Kvs[0].Value
	err := json.Unmarshal(podInstanceValue, &podInstance)
	if err != nil {
		fmt.Printf("%v\n", err)
		panic(err)
	}
	etcd.Delete(cli, podInstanceKey)

	//更新PodInstanceIDList
	podInstanceIDList := make([]string, 0)
	tmpList := make([]string, 0)
	kvs := etcd.Get(cli, def.PodInstanceListID).Kvs
	if len(kvs) != 0 {
		podInstanceIDListValue := kvs[0].Value
		err := json.Unmarshal(podInstanceIDListValue, &podInstanceIDList)
		if err != nil {
			fmt.Printf("%v\n", err)
			panic(err)
		}
	}
	for _, podInstanceID := range podInstanceIDList {
		if podInstanceID != podInstance.ID {
			tmpList = append(tmpList, podInstanceID)
		}
	}
	podInstanceIDList = tmpList
	podInstanceIDListValue, err := json.Marshal(podInstanceIDList)
	if err != nil {
		fmt.Printf("%v\n", err)
		panic(err)
	}
	etcd.Put(cli, def.PodInstanceListID, string(podInstanceIDListValue))

	//更新相应node中的PodInstances列表
	//nodeKey := "/node/" + strconv.Itoa(int(podInstance.NodeID))
	//nodeValue := etcd.Get(cli, nodeKey).Kvs[0].Value
	//var node def.Node
	//err = json.Unmarshal(nodeValue, &node)
	//if err != nil {
	//	fmt.Printf("%v\n", err)
	//	panic(err)
	//}
	//podInstances := make([]*def.PodInstance, len(node.PodInstances)-1)
	//podInstanceList := make([]string, len(node.PodInstances)-1)
	//for _, pi := range node.PodInstances {
	//	if pi.Pod.Metadata.Name != podInstance.Pod.Metadata.Name {
	//		podInstances = append(podInstances, pi)
	//		podInstanceList = append(podInstanceList, "/nodePodInstance/"+pi.Metadata.Name)
	//	}
	//}
	//node.PodInstances = podInstances
	//nodeValue, err = json.Marshal(node)
	//if err != nil {
	//	fmt.Printf("%v\n", err)
	//	panic(err)
	//}
	//etcd.Put(cli, nodeKey, string(nodeValue))

	//更新kubelet watch的node-PodInstance table
	//nodePIKey := "/nodePodInstance/" + strconv.Itoa(int(podInstance.NodeID))
	//nodePIValue, err := json.Marshal(podInstanceList)
	//if err != nil {
	//	fmt.Printf("%v\n", err)
	//	panic(err)
	//}
	//etcd.Put(cli, nodePIKey, string(nodePIValue))

	return true, podInstance
}

// CheckDeleteInService 检查是否需要将新删除的pod从之前创建好的service当中删除
func CheckDeleteInService(cli *clientv3.Client, podInstance def.PodInstance) []def.Service {
	servicePrefix := "/service/"
	kvs := etcd.GetWithPrefix(cli, servicePrefix).Kvs
	service := def.Service{}
	serviceValue := make([]byte, 0)
	serviceList := make([]def.Service, 0)
	for _, kv := range kvs {
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
			service.PortsBindings = RemovePodInstanceFromService(podInstance, service)
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

func RemovePodInstanceFromService(podInstance def.PodInstance, service def.Service) []def.PortsBindings {
	tmpBindings := make([]def.PortsBindings, 0)
	for _, binding := range service.PortsBindings {
		tmpBinding := def.PortsBindings{
			Ports: binding.Ports,
		}
		tmpEndpoints := make([]string, 0)
		for _, tmpEndpoint := range binding.Endpoints {
			tmpEndpointIP := (strings.Split(tmpEndpoint, ":"))[0]
			fmt.Println("tmpEndpointIP:   ", tmpEndpointIP)
			fmt.Println("podInstance.IP:   ", podInstance.IP)
			if tmpEndpointIP != podInstance.IP {
				tmpEndpoints = append(tmpEndpoints, tmpEndpointIP)
			}
		}
		tmpBinding.Endpoints = tmpEndpoints
		tmpBindings = append(tmpBindings, tmpBinding)
	}
	fmt.Println("tmpBindings:  ", tmpBindings)

	return tmpBindings
}
