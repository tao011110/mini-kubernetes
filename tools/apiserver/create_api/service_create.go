package create_api

import (
	"encoding/json"
	"fmt"
	"mini-kubernetes/tools/def"
	"mini-kubernetes/tools/etcd"
	"strconv"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

func CreateClusterIPService(cli *clientv3.Client, service_c def.ClusterIPSvc) def.Service {
	service := def.Service{
		Name:      service_c.Metadata.Name,
		Selector:  service_c.Spec.Selector,
		Type:      service_c.Spec.Type,
		ClusterIP: service_c.Spec.ClusterIP,
		StartTime: time.Now(),
	}
	portsBindingsList := make([]def.PortsBindings, 0)
	fmt.Printf("service.ClusterIP is %s", service.ClusterIP)

	podInstancePrefix := "/podInstance/"
	kvs := etcd.GetWithPrefix(cli, podInstancePrefix).Kvs
	podInstance := def.PodInstance{}
	podInstanceValue := make([]byte, 0)
	for _, ports := range service_c.Spec.Ports {
		endpoints := make([]string, 0)
		for _, kv := range kvs {
			podInstanceValue = kv.Value
			err := json.Unmarshal(podInstanceValue, &podInstance)
			if err != nil {
				fmt.Printf("%v\n", err)
				panic(err)
			}
			if podInstance.Pod.Metadata.Label == service.Selector.Name {
				for _, container := range podInstance.Pod.Spec.Containers {
					for _, portMapping := range container.PortMappings {
						containerPort := strconv.Itoa(int(portMapping.ContainerPort))
						fmt.Println(ports.TargetPort)
						fmt.Println(containerPort)
						if ports.TargetPort == containerPort {
							endpoints = append(endpoints, podInstance.IP+":"+ports.TargetPort)
							fmt.Printf("podInstance.IP is %v\n", podInstance.IP)
						}
					}
				}
			}
		}
		fmt.Printf("Endpoints is %v\n", endpoints)
		portsBindings := def.PortsBindings{
			Ports:     ports,
			Endpoints: endpoints,
		}
		portsBindingsList = append(portsBindingsList, portsBindings)
	}
	service.PortsBindings = portsBindingsList

	// write into etcd
	serviceKey := "/service/" + service.Name
	serviceValue, err := json.Marshal(service)
	if err != nil {
		fmt.Printf("%v\n", err)
		panic(err)
	}
	etcd.Put(cli, serviceKey, string(serviceValue))

	return service
}

func CreateNodePortService(cli *clientv3.Client, service_n def.NodePortSvc) def.Service {
	service := def.Service{
		Name:      service_n.Metadata.Name,
		Selector:  service_n.Spec.Selector,
		Type:      service_n.Spec.Type,
		ClusterIP: service_n.Spec.ClusterIP,
	}
	portsBindingsList := make([]def.PortsBindings, 0)
	fmt.Printf("service.ClusterIP is %s\n", service.ClusterIP)

	podInstancePrefix := "/podInstance/"
	kvs := etcd.GetWithPrefix(cli, podInstancePrefix).Kvs
	podInstance := def.PodInstance{}
	podInstanceValue := make([]byte, 0)
	fmt.Printf("service_c.Spec.Ports are %v\n", service_n.Spec.Ports)
	for _, ports := range service_n.Spec.Ports {
		fmt.Println(ports)
		endpoints := make([]string, 0)
		for _, kv := range kvs {
			podInstanceValue = kv.Value
			err := json.Unmarshal(podInstanceValue, &podInstance)
			if err != nil {
				fmt.Printf("%v\n", err)
				panic(err)
			}
			fmt.Printf("podInstance.Pod.Metadata.Label is %s\n", podInstance.Pod.Metadata.Label)
			fmt.Printf("service.Selector.Name is %s\n", service.Selector.Name)
			if podInstance.Pod.Metadata.Label == service.Selector.Name {
				for _, container := range podInstance.Pod.Spec.Containers {
					for _, portMapping := range container.PortMappings {
						containerPort := strconv.Itoa(int(portMapping.ContainerPort))
						fmt.Println(ports.TargetPort)
						fmt.Println(containerPort)
						if ports.TargetPort == containerPort {
							endpoints = append(endpoints, podInstance.IP+":"+ports.TargetPort)
							fmt.Printf("podInstance.IP is %v\n", podInstance.IP)
						}
					}
				}
			}
		}
		fmt.Printf("Endpoints is %v\n", endpoints)
		portsBindings := def.PortsBindings{
			Ports:     ports,
			Endpoints: endpoints,
		}
		portsBindingsList = append(portsBindingsList, portsBindings)
	}
	service.PortsBindings = portsBindingsList

	// write into etcd
	serviceKey := "/service/" + service.Name
	serviceValue, err := json.Marshal(service)
	if err != nil {
		fmt.Printf("%v\n", err)
		panic(err)
	}
	etcd.Put(cli, serviceKey, string(serviceValue))

	return service
}
