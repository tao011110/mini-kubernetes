package yaml_test

import (
	"encoding/json"
	"fmt"
	"github.com/YOUR-USER-OR-ORG-NAME/YOUR-REPO-NAME/tools/yaml"
	"log"
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
