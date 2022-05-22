package main

import (
	"fmt"
	"os"
	"os/exec"
	//"testing"
)

func main() {
	//var cmdList []string
	//cmdList = append(cmdList, `echo -e \"server:\\n  port: 80\\nspring:\\n  application:\\n    name: zuul\\nzuul:\\n  routes:\\n    route0:\\n      path: testPath/**\\n      url: 1.1.1.1:8080\\n\" > /home/zuul/src/main/resources/application.yaml`)
	//image_factory.ImageFactory(def.GatewayImage, "gatewaytest2", cmdList)
	file, _ := os.OpenFile("test.sh", os.O_RDWR, os.ModeAppend)
	cmd := exec.Command("docker", "exec", "test", "/bin/bash", "-c", fmt.Sprintf("'%s'", "command"))
	_, _ = file.Write([]byte(cmd.String()))
}
