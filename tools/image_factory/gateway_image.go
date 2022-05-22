package image_factory

import (
	"fmt"
	"mini-kubernetes/tools/def"
	"mini-kubernetes/tools/gateway"
)

func MakeGatewayImage(dns *def.DNSDetail, nameGatewayImageName string) {
	fileStr := gateway.GenerateApplicationYaml(*dns)
	echoCmd := EchoFactory(fileStr, def.GatewayRoutesConfigPathInImage)
	fmt.Println("echoCmd")
	fmt.Println(echoCmd)
	cmds := []string{echoCmd, def.GatewayPackageCmd}
	ImageFactory(def.GatewayImage, nameGatewayImageName, cmds)
}
