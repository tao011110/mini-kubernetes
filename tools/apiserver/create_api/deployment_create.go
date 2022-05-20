package create_api

import (
	"encoding/json"
	"fmt"
	"github.com/jakehl/goid"
	clientv3 "go.etcd.io/etcd/client/v3"
	"mini-kubernetes/tools/def"
	"mini-kubernetes/tools/etcd"
	"time"
)

func CreateDeployment(cli *clientv3.Client, deployment def.Deployment) {
	// Put deployment's name into deployment_list_name
	{
		deploymentListNameKey := def.DeploymentListName
		deploymentListNameValue := etcd.Get(cli, deploymentListNameKey).Kvs[0].Value
		list := make([]string, 0)
		err := json.Unmarshal(deploymentListNameValue, &list)
		if err != nil {
			fmt.Printf("%v\n", err)
			panic(err)
		}
		list = append(list, deployment.Metadata.Name)
		deploymentListNameValue, err = json.Marshal(list)
		if err != nil {
			fmt.Printf("%v\n", err)
			panic(err)
		}
		etcd.Put(cli, deploymentListNameKey, string(deploymentListNameValue))
	}

	podName := deployment.Metadata.Name + "-pod-" + goid.NewV4UUID().String()
	// Parse Deployment' meta into ParsedDeployment, and put it into etcd
	{
		parsedDeployment := def.ParsedDeployment{
			Name:        deployment.Metadata.Name,
			ReplicasNum: int(deployment.Spec.Replicas),
			PodName:     podName,
			StartTime:   time.Now(),
		}
		parsedDeploymentKey := def.GetKeyOfDeployment(parsedDeployment.Name)
		parsedDeploymentValue, err := json.Marshal(parsedDeployment)
		if err != nil {
			fmt.Printf("%v\n", err)
			panic(err)
		}
		etcd.Put(cli, parsedDeploymentKey, string(parsedDeploymentValue))
	}

	// Put ParsedDeployment's pod into etcd
	{
		containers := make([]def.Container, 0)
		volumes := make([]def.Volume, 0)
		for _, container := range deployment.Spec.Template.Spec.Containers {
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

		for _, volume := range deployment.Spec.Template.Spec.Volumes {
			tmp := def.Volume{
				Name:     volume.Name,
				HostPath: volume.HostPath.Path,
			}
			volumes = append(volumes, tmp)
		}

		pod := def.Pod{
			ApiVersion: deployment.ApiVersion,
			Kind:       "Pod",
			Metadata: def.PodMeta{
				Name:   podName,
				Labels: def.PodLabels(deployment.Spec.Template.Metadata.Labels),
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
