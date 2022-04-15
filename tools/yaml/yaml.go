package yaml

import (
	"os"
	"log"
	"fmt"
	"encoding/json"
	"gopkg.in/yaml.v2"
)

type Volume struct {
	Name string `yaml:"name"`
	HostPath string `yaml:"hostPath"`
}

type VolumeMount struct {
	Name string `yaml:"name"`
	MountPath string  `yaml:"mountPath"`
}

type PortMapping struct {
	Name string `yaml:"name"`
	ContainerPort uint16 `yaml:"containerPort"`
	HostPort uint16 `yaml:"hostPort"`
	Protocol string `yaml:"protocol"`
}

type Resource struct {
	ResourceLimit Limit `yaml:"limits"`
	ResourceRequest Request `yaml:"requests"`
}

type Limit struct {
	CPU string `yaml:"cpu"`
	Memory string `yaml:"memory"`
}

type Request struct {
	CPU string `yaml:"cpu"`
	Memory string `yaml:"memory"`
}

type Container struct {
	Name string `yaml:"name"`
	Image string `yaml:"image"`
	Commands []string `yaml:"command"`
	Args []string `yaml:"args"`
	WorkingDir string `yaml:"workingDir"`
	VolumeMounts []VolumeMount `yaml:"volumeMounts"`
	PortMappings []PortMapping `yaml:"ports"`
	Resources Resource `yaml:"resources"`
}

type HttpHeaderPair struct {
	Name string `yaml:"name"`
	Value string `yaml:"value"`
}

type HttpRequest struct {
	Scheme string `yaml:"scheme"`
	HttpHeaders []HttpHeaderPair `yaml:"HttpHeaders"`
	Path string `yaml:"path"`
	Port uint16 `yaml:"port"`
}

type Exec struct {
	Command string `yaml:"command"`
}

type LivenessProbe struct {
	InitialDelaySeconds  uint32 `yaml:"initialDelaySeconds"`
	TimeoutSeconds       uint32 `yaml:"timeoutSeconds"`
	PeriodSeconds        uint32 `yaml:"periodSeconds"`
	FailureThreshold     uint32 `yaml:"failureThreshold"`
	SuccessThreshold     uint32 `yaml:"successThreshold"`

	Exec Exec `yaml:"exec"`
	HttpGetRequest HttpRequest `yaml:"httpGet"`
}

type Spec struct {
	Containers []Container	`yaml:"containers"`
	LivenessProbe LivenessProbe `yaml:"livenessProbe"`
	Volumes []Volume `yaml:"volumes"`
}

type Meta struct {
	Name string `yaml:"name"`
}

type Pod struct {
	ApiVersion string `yaml:"apiVersion"`
	Kind string `yaml:"kind"`
	Metadata Meta `yaml:"metadata"`
	// NodeID和NodeSelector迭代三再加
	Spec Spec `yaml:"spec"`
}


func ReadYamlConfig(path string) (*Pod,error){
	pod := &Pod{}
	if f, err := os.Open(path); err != nil {
		return nil,err
	} else {
		err := yaml.NewDecoder(f).Decode(pod)
		if err != nil {
			return nil,err
		}
		if (*pod).ApiVersion != "v1" {
			fmt.Println("apiVersion should be v1!")
			return nil,err
		} else if(*pod).Kind != "Pod" {
			fmt.Println("kind should be Pod!")
			return nil,err
		}
	}
	fmt.Println("pod: ", pod)
	return pod,nil
}
/*
func main() {
	pod,err := ReadYamlConfig("tmp.yaml")
	if err != nil {
		log.Fatal(err)
	}

	byts, err := json.Marshal(pod)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("string:",string(byts))
}
*/
