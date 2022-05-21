package yaml_test

import (
	"encoding/json"
	"fmt"
	"log"
	"mini-kubernetes/tools/yaml"
	"testing"
)

func Test(t *testing.T) {
	// test read pod yaml
	pod, err := yaml.ReadPodYamlConfig("tmp.yaml")
	if err != nil {
		log.Fatal(err)
	}
	byts, err := json.Marshal(pod)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("string:", string(byts))

	// test read deployment yaml
	dep, err := yaml.ReadDeploymentConfig("dep.yaml")
	if err != nil {
		log.Fatal(err)
	}
	byts, err = json.Marshal(dep)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("string:", string(byts))

	// test read autoscaler yaml
	auto, err := yaml.ReadAutoScalerConfig("auto.yaml")
	if err != nil {
		log.Fatal(err)
	}

	byts, err = json.Marshal(auto)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("string:", string(byts))

	// judge config type
	num, err := yaml.ReadConfig("auto.yaml")
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("typeid: %d\n", num)
	}

	num, err = yaml.ReadConfig("cluster.yaml")
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("typeid: %d\n", num)
	}
}
