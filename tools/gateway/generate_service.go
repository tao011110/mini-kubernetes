package gateway

import (
	"fmt"
	"mini-kubernetes/tools/def"
)

func generatePodPortMappings(dns *def.DNSDetail) []def.PortMapping {
	var mappings []def.PortMapping
	mappings = append(mappings, def.PortMapping{
		Name:          fmt.Sprintf("gateway_port_%d", 80),
		ContainerPort: 80,
		//HostPort:      80,
		Protocol: "TCP",
	})
	return mappings
}

func generateServicePorts(dns *def.DNSDetail) []def.PortPair {
	var mappings []def.PortPair
	for _, path := range dns.Paths {
		mappings = append(mappings, def.PortPair{
			Port:       path.Port,
			TargetPort: fmt.Sprintf(`%d`, path.Port),
			Protocol:   "TCP",
		})
	}
	return mappings
}

// GenerateGatewayPod TODO: 创建的容器并不能直接start?
func GenerateGatewayPod(dns def.DNSDetail, imageName string) (pod def.Pod) {
	//injectionRoutesCommand := fmt.Sprintf(
	//	"echo -e \"%s\" > %s",
	//	GenerateApplicationYaml(dns),
	//	def.GatewayRoutesConfigPathInImage)
	//fmt.Println("injectionRoutesCommand")
	//fmt.Println(injectionRoutesCommand)
	//packageAndRunCommand := fmt.Sprintf(
	//	"./%s",
	//	def.GatewayPackageAndRunScriptPath)
	//fmt.Println("packageAndRunCommand")
	//fmt.Println(packageAndRunCommand)
	gatewayContainerResource := def.Resource{
		ResourceLimit: def.Limit{
			CPU:    `3`,
			Memory: `3G`,
		},
		ResourceRequest: def.Request{
			CPU:    `2`,
			Memory: `2.5G`,
		},
	}
	containerName := fmt.Sprintf("gateway_container_%s_name", dns.Name)
	podName := fmt.Sprintf("gateway_pod_%s_name", dns.Name)
	podLabel := fmt.Sprintf("gateway_pod_%s_label", dns.Name)
	//serviceName := fmt.Sprintf("gateway_service_%s_name", dns.Name)

	pod = def.Pod{
		ApiVersion: `v1`,
		Kind:       `Pod`,
		Metadata: def.PodMeta{
			Name:  podName,
			Label: podLabel,
		},
		Spec: def.PodSpec{
			Containers: []def.Container{
				{
					Name:  containerName,
					Image: imageName,
					//Commands:     []string{injectionRoutesCommand, packageAndRunCommand},
					//Commands:     []string{packageAndRunCommand},
					//Args:         []string{},
					//WorkingDir:   "",
					//VolumeMounts: []def.VolumeMount{},
					PortMappings: generatePodPortMappings(&dns),
					Resources:    gatewayContainerResource,
				},
			},
			Volumes: []def.Volume{},
		},
	}

	//service = def.ClusterIPSvc{
	//	ApiVersion: `v1`,
	//	Kind:       `Service`,
	//	Metadata: def.Meta{
	//		Name: serviceName,
	//	},
	//	Spec: def.Spec{
	//		Type:  `ClusterIP`,
	//		Ports: generateServicePorts(&dns),
	//		Selector: def.Selector{
	//			Name: podLabel, /*TODO: maybe wrong*/
	//		},
	//	},
	//}
	return
}
