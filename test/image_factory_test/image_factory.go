package main

import (
	"mini-kubernetes/tools/def"
	"mini-kubernetes/tools/image_factory"
	//"testing"
)

func main() {
	var cmdList []string
	cmdList = append(cmdList, `echo -e \"server:\\n  port: 80\\nspring:\\n  application:\\n    name: zuul\\nzuul:\\n  routes:\\n    route0:\\n      path: testPath/**\\n      url: 1.1.1.1:8080\\n\" > /home/zuul/src/main/resources/application.yaml`)
	image_factory.ImageFactory(def.GatewayImage, "gatewaytest2", cmdList)
}
