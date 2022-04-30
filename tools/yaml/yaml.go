package yaml

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"mini-kubernetes/tools/pod"
	"os"
)

func ReadYamlConfig(path string) (*pod.Pod, error) {
	pod := &pod.Pod{}
	if f, err := os.Open(path); err != nil {
		return nil, err
	} else {
		err := yaml.NewDecoder(f).Decode(pod)
		if err != nil {
			return nil, err
		}
		if (*pod).ApiVersion != "v1" {
			fmt.Println("apiVersion should be v1!")
			return nil, err
		} else if (*pod).Kind != "Pod" {
			fmt.Println("kind should be Pod!")
			return nil, err
		}
	}
	fmt.Println("pod: ", pod)
	return pod, nil
}
