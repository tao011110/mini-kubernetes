package create_api

import (
	"encoding/json"
	"fmt"
	"mini-kubernetes/tools/def"
	"mini-kubernetes/tools/etcd"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

func CreateAutoscaler(cli *clientv3.Client, autoscaler def.Autoscaler) {
	podName := def.GetPodNameOfAutoscaler(autoscaler.Metadata.Name)
	// Parse autoscaler' meta into ParsedDeployment, and put it into etcd
	{
		//CPUMaxValue, _ := strconv.ParseFloat(autoscaler.Spec.Metrics.CPU.TargetMaxValue, 64)
		//CPUMinValue, _ := strconv.ParseFloat(autoscaler.Spec.Metrics.CPU.TargetMinValue, 64)
		//memoryMaxValue, _ := strconv.ParseInt(autoscaler.Spec.Metrics.Memory.TargetMaxValue, 10, 64)
		//memoryMinValue, _ := strconv.ParseInt(autoscaler.Spec.Metrics.Memory.TargetMinValue, 10, 64)
		parsedAutoscaler := def.ParsedHorizontalPodAutoscaler{
			Name:           autoscaler.Metadata.Name,
			CPUMaxValue:    autoscaler.Spec.Metrics.CPU.TargetMaxValue,
			CPUMinValue:    autoscaler.Spec.Metrics.CPU.TargetMinValue,
			MemoryMaxValue: autoscaler.Spec.Metrics.Memory.TargetMaxValue,
			MemoryMinValue: autoscaler.Spec.Metrics.Memory.TargetMinValue,
			MaxReplicas:    int(autoscaler.Spec.MaxReplicas),
			MinReplicas:    int(autoscaler.Spec.MinReplicas),
			PodName:        podName,
			StartTime:      time.Now(),
		}
		autoscalerKey := def.GetKeyOfAutoscaler(autoscaler.Metadata.Name)
		autoscalerValue, err := json.Marshal(parsedAutoscaler)
		if err != nil {
			fmt.Printf("%v\n", err)
			panic(err)
		}
		etcd.Put(cli, autoscalerKey, string(autoscalerValue))
	}

	// Put autoscaler's pod into etcd
	{
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
				Name:  podName,
				Label: autoscaler.Spec.Template.Metadata.Labels.Name,
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

	// Put autoscaler's name into autoscaler_list_name
	{
		autoscalerListNameKey := def.HorizontalPodAutoscalerListName
		kvs := etcd.Get(cli, autoscalerListNameKey).Kvs
		list := make([]string, 0)
		if len(kvs) != 0 {
			autoscalerListNameValue := kvs[0].Value
			err := json.Unmarshal(autoscalerListNameValue, &list)
			if err != nil {
				fmt.Printf("%v\n", err)
				panic(err)
			}
		}
		list = append(list, autoscaler.Metadata.Name)
		autoscalerListNameValue, err := json.Marshal(list)
		if err != nil {
			fmt.Printf("%v\n", err)
			panic(err)
		}
		etcd.Put(cli, autoscalerListNameKey, string(autoscalerListNameValue))
	}
}
