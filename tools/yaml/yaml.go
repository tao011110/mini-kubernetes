package yaml

import (
	"fmt"
	"gopkg.in/yaml.v2"
	def "mini-kubernetes/tools/def"
	"os"
)

func ReadYamlConfig(path string) (*def.Pod, error) {
	pod := &def.Pod{}
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
