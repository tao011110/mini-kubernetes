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

func GenerateGatewayPod(dns def.DNSDetail, imageName string) (pod def.Pod) {
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
					Image:        imageName,
					PortMappings: generatePodPortMappings(&dns),
					Resources:    gatewayContainerResource,
					Commands:     []string{def.StartBash},
					Args:         []string{def.GatewayStartArgs},
				},
			},
			Volumes: []def.Volume{},
		},
	}
	return
}
