package create_api

import (
	"encoding/json"
	"fmt"
	"github.com/jakehl/goid"
	clientv3 "go.etcd.io/etcd/client/v3"
	"mini-kubernetes/tools/def"
	"mini-kubernetes/tools/etcd"
)

func CreateAutoscaler(cli *clientv3.Client, autoscaler def.Autoscaler) {
	// Put autoscaler's name into autoscaler_list_name
	{
		autoscalerListNameKey := def.HorizontalPodAutoscalerListName
		autoscalerListNameValue := etcd.Get(cli, autoscalerListNameKey).Kvs[0].Value
		list := make([]string, 0)
		err := json.Unmarshal(autoscalerListNameValue, &list)
		if err != nil {
			fmt.Printf("%v\n", err)
			panic(err)
		}
		list = append(list, autoscaler.Metadata.Name)
		autoscalerListNameValue, err = json.Marshal(list)
		if err != nil {
			fmt.Printf("%v\n", err)
			panic(err)
		}
		etcd.Put(cli, autoscalerListNameKey, string(autoscalerListNameValue))
	}

	// Parse autoscaler' meta into ParsedDeployment, and put it into etcd
	{
		autoscalerKey := def.GetKeyOfAutoscaler(autoscaler.Metadata.Name)
		autoscalerValue, err := json.Marshal(autoscaler)
		if err != nil {
			fmt.Printf("%v\n", err)
			panic(err)
		}
		etcd.Put(cli, autoscalerKey, string(autoscalerValue))
	}

	// Put autoscaler's pod into etcd
	{
		podName := autoscaler.Metadata.Name + "-pod-" + goid.NewV4UUID().String()
		containers := make([]def.Container, 0)
		volumes := make([]def.Volume, 0)
		for _, container := range autoscaler.Spec.Template.Spec.Containers {
			volumeMounts := make([]def.VolumeMount, 0)
			portMappings := make([]def.PortMapping, 0)
			for _, volumeMount := range container.VolumeMounts {
				vm := def.VolumeMount{
					Name:      volumeMount.Name,
					MountPath: volumeMount.MountPath,
				}
				volumeMounts = append(volumeMounts, vm)
			}
			for _, portMapping := range container.PortMappings {
				pm := def.PortMapping{
					Name:          portMapping.Name,
					ContainerPort: portMapping.ContainerPort,
					Protocol:      portMapping.Protocol,
				}
				portMappings = append(portMappings, pm)
			}
			tmp := def.Container{
				Name:         container.Name,
				Image:        container.Image,
				VolumeMounts: volumeMounts,
				PortMappings: portMappings,
			}
			containers = append(containers, tmp)
		}

		for _, volume := range autoscaler.Spec.Template.Spec.Volumes {
			tmp := def.Volume{
				Name:     volume.Name,
				HostPath: volume.HostPath.Path,
			}
			volumes = append(volumes, tmp)
		}

		pod := def.Pod{
			ApiVersion: autoscaler.ApiVersion,
			Kind:       "Pod",
			Metadata: def.PodMeta{
				Name:   podName,
				Labels: def.PodLabels(autoscaler.Spec.Template.Metadata.Labels),
			},
			Spec: def.PodSpec{
				Containers: containers,
				Volumes:    volumes,
			},
		}
		podKey := "/pod/" + pod.Metadata.Name
		podValue, err := json.Marshal(pod)
		if err != nil {
			fmt.Printf("%v\n", err)
			panic(err)
		}
		etcd.Put(cli, podKey, string(podValue))
	}
}
