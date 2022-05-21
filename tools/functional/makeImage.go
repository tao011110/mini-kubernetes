package functional

import (
	"fmt"
	"github.com/jakehl/goid"
	"mini-kubernetes/tools/def"
	"mini-kubernetes/tools/util"
	"strings"
)

func MakeFunctionalImage(function *def.Function) {
	pyString := util.ReadFile(function.Function)
	requirementsString := util.ReadFile(function.Requirements)
	pyString = strings.Replace(pyString, "\n", "\\n", -1)
	requirementsString = strings.Replace(requirementsString, "\n", "\\n", -1)
	cmdWritePy := fmt.Sprintf("echo -e \"%s\" > %s", pyString, def.PyHandlerPath)
	cmdWriteRequirements := fmt.Sprintf("echo -e \"%s\" > %s", requirementsString, def.RequirementsPath)
	cmdPrepare := fmt.Sprintf("./%s", def.PreparePath)
	imageName := fmt.Sprintf("image_%s_%d", function.Name, function.Version)
	// TODO: 拉取def.TemplateImage image, 启动, 初始命令为cmd_writePy, cmd_writeRequirements和cmdPrepare, limit适当(拉取依赖)
	// TODO: commit & push, image name为repository_name/imageName docker push repository_name/imageName
	function.Image = imageName
}

func GenerateFunctionPodAndService(function def.Function) (pod def.Pod, service def.ClusterIPSvc) {
	defaultContainerResource := def.Resource{
		ResourceLimit: def.Limit{
			CPU:    `3`,
			Memory: `3G`,
		},
		ResourceRequest: def.Request{
			CPU:    `1`,
			Memory: `500Mi`,
		},
	}
	id := fmt.Sprintf("functional_%s_%d_%s", function.Name, function.Version, goid.NewV4UUID().String())
	containerName := fmt.Sprintf("image_%s", id)
	podName := fmt.Sprintf("pod_%s", id)
	podLabel := fmt.Sprintf("pod_%s_label", id)
	serviceName := fmt.Sprintf("service_%s", id)
	cmdStart := fmt.Sprintf("./%s", def.StartPath)
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
					Image:        function.Image,
					Commands:     []string{cmdStart},
					Args:         []string{},
					WorkingDir:   "",
					VolumeMounts: []def.VolumeMount{},
					PortMappings: []def.PortMapping{
						{
							Name:          fmt.Sprintf("portmaping_%s", id),
							ContainerPort: 80,
							HostPort:      80,
							Protocol:      "HTTP",
						},
					},
					Resources: defaultContainerResource,
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
			Type: `ClusterIP`,
			Ports: []def.PortPair{{
				Port:       80,
				TargetPort: `80`,
				Protocol:   "HTTP",
			}},
			Selector: def.Selector{
				Name: podLabel, /*TODO: maybe wrong*/
			},
		},
	}
	return
}
