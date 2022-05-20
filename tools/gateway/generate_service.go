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
		HostPort:      80,
		Protocol:      "HTTP",
	})
	return mappings
}

func generateServicePorts(dns *def.DNSDetail) []def.PortPair {
	var mappings []def.PortPair
	for _, path := range dns.Paths {
		mappings = append(mappings, def.PortPair{
			Port:       path.Port,
			TargetPort: fmt.Sprintf(`%d`, path.Port),
			Protocol:   "HTTP",
		})
	}
	return mappings
}

func GenerateGatewayPodAndService(dns def.DNSDetail) (pod def.Pod, service def.ClusterIPSvc) {
	injectionRoutesCommand := fmt.Sprintf(
		"echo %s > %s",
		GenerateApplicationYaml(dns),
		def.GatewayRoutesConfigPathInImage)
	packageAndRunCommand := fmt.Sprintf(
		"./%s",
		def.GatewayPackageAndRunScriptPath)
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
	serviceName := fmt.Sprintf("gateway_service_%s_name", dns.Name)

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
					Name:         containerName,
					Image:        def.GatewayImage,
					Commands:     []string{injectionRoutesCommand, packageAndRunCommand},
					Args:         []string{},
					WorkingDir:   "",
					VolumeMounts: []def.VolumeMount{},
					PortMappings: generatePodPortMappings(&dns),
					Resources:    gatewayContainerResource,
				},
			},
			Volumes: []def.Volume{},
		},
	}

	service = def.ClusterIPSvc{
		ApiVersion: `v1`,
		Kind:       `Service`,
		Metadata: def.Meta{
			Name: serviceName,
		},
		Spec: def.Spec{
			Type:  `ClusterIP`,
			Ports: generateServicePorts(&dns),
			Selector: def.Selector{
				Name: podLabel, /*TODO: maybe wrong*/
			},
		},
	}
	return
}
