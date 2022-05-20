package def

import "time"

// ParsedDeployment
//下面不是deployment解析的直接结果
//为了复用pod相关的接口, 后请把template中的内容转为pod并为此pod分配全局唯一的name

type ParsedDeployment struct {
	Name        string    `json:"name"`
	ReplicasNum int       `json:"replicas_num"`
	PodName     string    `json:"pod_name"`
	StartTime   time.Time `json:"time"`
}

type ReplicasState struct {
	Desired     int `json:"Desired"`
	Updated     int `json:"Updated"`
	Total       int `json:"Total"`
	Available   int `json:"Available"`
	Unavailable int `json:"Unavailable"`
}

type DeploymentDetail struct {
	Name              string        `json:"name"`
	PodTemplate       Pod           `json:"podTemplate"`
	ReplicasState     ReplicasState `json:"replicasState"`
	CreationTimestamp time.Time     `json:"creationTimestamp"`
}

// DeploymentBrief 是针对 kubectl get deployment返回的信息
type DeploymentBrief struct {
	Name      string        `json:"name"`
	Ready     string        `json:"ready"`
	UpToDate  int           `json:"upToDate"`
	Available int           `json:"available"`
	Age       time.Duration `json:"age"`
}

// 以下是deployment解析的直接结果
type DepMeta struct {
	Name string `yaml:"name" json:"name"`
}

type TempLabels struct {
	Name string `yaml:"name" json:"name"`
}

type TempMeta struct {
	Labels TempLabels `yaml:"labels" json:"labels"`
}

type TempVolumeMount struct {
	Name      string `yaml:"name" json:"name"`
	MountPath string `yaml:"mountPath" json:"mount_path"`
}

type TempPortMapping struct {
	Name          string `yaml:"name" json:"name"`
	ContainerPort uint16 `yaml:"containerPort" json:"container_port"`
	Protocol      string `yaml:"protocol" json:"protocol"`
}

type TempContainer struct {
	Name         string            `yaml:"name" json:"name"`
	Image        string            `yaml:"image" json:"image"`
	VolumeMounts []TempVolumeMount `yaml:"volumeMounts" json:"volumeMounts"`
	PortMappings []TempPortMapping `yaml:"ports" json:"portMappings"`
}

type TempHostPath struct {
	Path string `yaml:"path" json:"path"`
}

type TempVolume struct {
	Name     string       `yaml:"name" json:"name"`
	HostPath TempHostPath `yaml:"hostPath" json:"host_path"`
}

type TempSpec struct {
	Containers []TempContainer `yaml:"containers" json:"containers"`
	Volumes    []TempVolume    `yaml:"volumes" json:"volumes"`
}

type DepTemp struct {
	Metadata TempMeta `yaml:"metadata" json:"metadata"`
	Spec     TempSpec `yaml:"spec" json:"spec"`
}

type DepSpec struct {
	Replicas uint64  `yaml:"replicas" json:"replicas"`
	Template DepTemp `yaml:"template" json:"template"`
}

type Deployment struct {
	ApiVersion string  `yaml:"apiVersion" json:"api_version"`
	Kind       string  `yaml:"kind" json:"kind"`
	Metadata   DepMeta `yaml:"metadata" json:"metadata"`
	Spec       DepSpec `yaml:"spec" json:"spec"`
}
