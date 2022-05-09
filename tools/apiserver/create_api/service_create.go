package create_api

import (
	"encoding/json"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"mini-kubernetes/tools/def"
	"mini-kubernetes/tools/etcd"
	"strconv"
)

func CreateClusterIPService(cli *clientv3.Client, service_c def.ClusterIP) def.Service {
	service := def.Service{
		Name:     service_c.Metadata.Name,
		Selector: service_c.Spec.Selector,
		Type:     service_c.Spec.Type,
		IP:       service_c.Spec.ClusterIP,
	}
	portsBindingsList := make([]def.PortsBindings, 0)
	fmt.Printf("service.IP is %s", service.IP)

	podInstanceKey := "/podInstance/"
	kvs := etcd.GetWithPrefix(cli, podInstanceKey).Kvs
	podInstance := def.PodInstance{}
	podInstanceValue := make([]byte, 0)
	if len(kvs) != 0 {
		for _, ports := range service_c.Spec.Ports {
			fmt.Println(ports)
			endpoints := make([]string, 0)
			for _, kv := range kvs {
				podInstanceValue = kv.Value
				err := json.Unmarshal(podInstanceValue, &podInstance)
				if err != nil {
					fmt.Printf("%v\n", err)
					panic(err)
				}
				if podInstance.Pod.Metadata.Labels.Name == service.Selector.Name {
					for _, container := range podInstance.Pod.Spec.Containers {
						for _, portMapping := range container.PortMappings {
							containerPort := strconv.Itoa(int(portMapping.ContainerPort))
							fmt.Println(ports.TargetPort)
							fmt.Println(containerPort)
							if ports.TargetPort == containerPort {
								endpoints = append(endpoints, podInstance.IP)
								fmt.Printf("podInstance.IP is %v\n", podInstance.IP)
							}
						}
					}
				}
			}
			if len(endpoints) != 0 {
				fmt.Printf("Endpoints is %v\n", endpoints)
				portsBindings := def.PortsBindings{
					Ports:     ports,
					Endpoints: endpoints,
				}
				portsBindingsList = append(portsBindingsList, portsBindings)
			}
		}
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
