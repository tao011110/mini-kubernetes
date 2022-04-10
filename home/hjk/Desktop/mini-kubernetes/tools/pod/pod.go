// created by hjk 2022.4.10

package pod

type Volume struct {
	Name string
	HostPath string
}

type VolumeMount struct {
	Name string
	MountPath string
}

type PortMapping struct {
	Name string
	ContainerPort uint16
	HostPort uint16
	Protocol string
}

type Resource struct {
	CPU uint16
	Memory uint64
}

type Container struct {
	Name string
	Image string
	Commands []string
	Args []string
	WorkingDir string
	VolumeMounts []VolumeMount
	PortMappings []PortMapping
	ResourceLimit Resource
	ResourceRequest Resource
}

type HttpHeaderPair struct {
	Name string
	Value string
}

type HttpRequest struct {
	Scheme string
	HttpHeaders []HttpHeaderPair
	Path string
	Port uint16
}

type Command string

type LivenessProbe struct {
	InitialDelaySeconds  uint32
	TimeoutSeconds       uint32
	PeriodSeconds        uint32
	FailureThreshold     uint32
	SuccessThreshold     uint32

	HttpHeaders []HttpHeaderPair
	Command Command
	HttpGetRequest HttpRequest
}

type Pod struct {
	Name string
	Containers []Container
	Volumes []Volume
	LivenessProbe LivenessProbe
	NodeName string
	//nodeSelector 迭代三再加
}
