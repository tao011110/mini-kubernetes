package yaml_test

import (
	"encoding/json"
	"fmt"
	"log"
	"mini-kubernetes/tools/yaml"
	"testing"
)

func Test(t *testing.T) {
	pod, err := yaml.ReadYamlConfig("tmp.yaml")
	if err != nil {
		log.Fatal(err)
	}

	byts, err := json.Marshal(pod)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("string:", string(byts))
}
