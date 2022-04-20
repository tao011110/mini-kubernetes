package def

import (
	"github.com/golang/protobuf/ptypes/timestamp"
	"net"
)

type Volume struct {
	Name     string `yaml:"name" json:"name"`
	HostPath string `yaml:"hostPath" json:"host_path"`
}

type VolumeMount struct {
	Name      string `yaml:"name" json:"name"`
	MountPath string `yaml:"mountPath" json:"mount_path"`
}

type PortMapping struct {
	Name          string `yaml:"name" json:"name"`
	ContainerPort uint16 `yaml:"containerPort" json:"container_port"`
	HostPort      uint16 `yaml:"hostPort" json:"host_port"`
	Protocol      string `yaml:"protocol" json:"protocol"`
}

type Resource struct {
	ResourceLimit   Limit   `yaml:"limits" json:"resource_limit"`
	ResourceRequest Request `yaml:"requests" json:"resource_request"`
}

type Limit struct {
	CPU    string `yaml:"cpu" json:"cpu"`
	Memory string `yaml:"memory" json:"memory"`
}

type Request struct {
	CPU    string `yaml:"cpu" json:"cpu"`
	Memory string `yaml:"memory" json:"memory"`
}

type Container struct {
	Name         string        `yaml:"name" json:"name"`
	Image        string        `yaml:"image" json:"image"`
	Commands     []string      `yaml:"command" json:"commands"`
	Args         []string      `yaml:"args" json:"args"`
	WorkingDir   string        `yaml:"workingDir" json:"workingDir"`
	VolumeMounts []VolumeMount `yaml:"volumeMounts" json:"volumeMounts"`
	PortMappings []PortMapping `yaml:"ports" json:"portMappings"`
	Resources    Resource      `yaml:"resources" json:"resources"`
}

type HttpHeaderPair struct {
	Name  string `yaml:"name" json:"name"`
	Value string `yaml:"value" json:"value"`
}

type HttpRequest struct {
	Scheme      string           `yaml:"scheme" json:"scheme"`
	HttpHeaders []HttpHeaderPair `yaml:"HttpHeaders" json:"httpHeaders"`
	Path        string           `yaml:"path" json:"path"`
	Port        uint16           `yaml:"port" json:"port"`
}

type Exec struct {
	Command string `yaml:"command" json:"command"`
}

type LivenessProbe struct {
	InitialDelaySeconds uint32 `yaml:"initialDelaySeconds" json:"initial_delay_seconds"`
	TimeoutSeconds      uint32 `yaml:"timeoutSeconds" json:"timeout_seconds"`
	PeriodSeconds       uint32 `yaml:"periodSeconds" json:"period_seconds"`
	FailureThreshold    uint32 `yaml:"failureThreshold" json:"failure_threshold"`
	SuccessThreshold    uint32 `yaml:"successThreshold" json:"success_threshold"`

	Exec           Exec        `yaml:"exec" json:"exec"`
	HttpGetRequest HttpRequest `yaml:"httpGet" json:"http_get_request"`
}

type Spec struct {
	Containers    []Container   `yaml:"containers" json:"containers"`
	LivenessProbe LivenessProbe `yaml:"livenessProbe" json:"livenessProbe"`
	Volumes       []Volume      `yaml:"volumes" json:"volumes"`
}

type Meta struct {
	Name string `yaml:"name" json:"name"`
}

type Pod struct {
	ApiVersion string `yaml:"apiVersion" json:"api_version"`
	Kind       string `yaml:"kind" json:"kind"`
	Metadata   Meta   `yaml:"metadata" json:"metadata"`
	// NodeID和NodeSelector迭代三再加
	Spec Spec `yaml:"spec" json:"spec"`
}

const (
	PENDING   uint8 = 0 //docker并未运行, 如正在下载镜像等
	RUNNING   uint8 = 1 //
	SUCCEEDED uint8 = 2 //所有docker均正常结束
	FAILED    uint8 = 3
	UNKNOWN   uint8 = 4 //连接不到node, 用于master
)

type PodInstance struct {
	Pod
	IP              net.IPAddr          `json:"ip"`
	NodeID          uint64              `json:"nodeID"`
	StartTime       timestamp.Timestamp `json:"startTime"`
	Status          uint8               `json:"status"`
	ContainerStatus []uint8             `json:"containerStatus"`
	RestartCount    uint64              `json:"restartCount"`
}
