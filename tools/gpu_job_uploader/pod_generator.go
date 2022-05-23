package gpu_job_uploader

import (
	"fmt"
	"mini-kubernetes/tools/def"
)

func GenerateGpuJobUploaderPod(job *def.GPUJob) def.Pod {
	generateImage(job)
	defualtResource := def.Resource{
		ResourceLimit: def.Limit{
			CPU:    `1`,
			Memory: `500M`,
		},
		ResourceRequest: def.Request{
			CPU:    `1`,
			Memory: `500M`,
		},
	}
	containerName := fmt.Sprintf("gpuUploader_container_%s_name", job.Name)
	podName := fmt.Sprintf("gpuUploader_pod_%s_name", job.Name)
	podLabel := fmt.Sprintf("gpuUploader_pod_%s_label", job.Name)

	return def.Pod{
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
					Image: job.ImageName,
					PortMappings: []def.PortMapping{{
						Name:          "port_mapping_80",
						ContainerPort: 80,
						HostPort:      80,
						Protocol:      "HTTP",
					}},
					Resources: defualtResource,
					Commands:  []string{def.GPUJobUploaderRunCmd},
				},
			},
			Volumes: []def.Volume{},
		},
	}
}
