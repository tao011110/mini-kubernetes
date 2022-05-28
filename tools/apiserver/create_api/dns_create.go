package create_api

import (
	"encoding/json"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"mini-kubernetes/tools/apiserver/get_api"
	"mini-kubernetes/tools/def"
	"mini-kubernetes/tools/etcd"
	"mini-kubernetes/tools/gateway"
	"mini-kubernetes/tools/image_factory"
)

func CreateDNS(cli *clientv3.Client, dns def.DNS) (def.DNSDetail, string) {
	// Generate DNSDetail from DNS
	// Create pathPairDetails(containing clusterIP service) from DNS's path
	pathPairDetails := make([]def.PathPairDetail, 0)
	for _, pair := range dns.Paths {
		svc, _ := get_api.GetService(cli, pair.Service)
		ports := make([]def.PortPair, 0)
		for _, binds := range svc.PortsBindings {
			ports = append(ports, binds.Ports)
		}
		clusterIPSvc := def.ClusterIPSvc{
			Metadata: def.Meta{
				Name: svc.Name,
			},
			Spec: def.Spec{
				ClusterIP: svc.ClusterIP,
				Ports:     ports,
				Selector:  svc.Selector,
			},
		}
		pathPairDetail := def.PathPairDetail{
			Path:    pair.Path,
			Port:    pair.Port,
			Service: clusterIPSvc,
		}
		fmt.Printf("pathPairDetail is %v\n", pathPairDetail)
		pathPairDetails = append(pathPairDetails, pathPairDetail)
	}
	fmt.Printf("pathPairDetails is %v\n", pathPairDetails)

	dnsDetailKey := "/DNS/" + dns.Name
	dnsDetail := def.DNSDetail{
		Kind:  dns.Kind,
		Name:  dns.Name,
		Host:  dns.Host,
		Paths: pathPairDetails,
	}
	dnsDetailValue, err := json.Marshal(dnsDetail)
	if err != nil {
		fmt.Printf("%v\n", err)
		panic(err)
	}
	etcd.Put(cli, dnsDetailKey, string(dnsDetailValue))
	fmt.Printf("dnsDetail is  %v\n", dnsDetail)

	imageName := "tmpForGateway"
	image_factory.MakeGatewayImage(&dnsDetail, imageName)

	// Create Gateway, then put its pod and service into etcd
	gatewayPod := gateway.GenerateGatewayPod(dnsDetail, def.RgistryAddr+imageName)
	gatewayPod.Spec.Containers[0].Image = def.RgistryAddr + imageName
	podInstance := CreatePod(cli, gatewayPod)

	return dnsDetail, podInstance.ID
}
