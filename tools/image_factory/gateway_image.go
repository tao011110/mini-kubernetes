package image_factory

import (
	"mini-kubernetes/tools/def"
	"mini-kubernetes/tools/gateway"
)

func MakeGatewayImage(dns *def.DNSDetail, nameGatewayImageName string) {
	fileStr := gateway.GenerateApplicationYaml(*dns)
	echoCmd := EchoFactory(fileStr, def.GatewayRoutesConfigPathInImage)
	cmds := []string{echoCmd, def.GatewayPackageCmd}
	ImageFactory(def.GatewayImage, nameGatewayImageName, cmds)
}
