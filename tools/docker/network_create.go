package docker

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
)

func CreateNetBridge(cniIP string) string {
	name := "miniK8S-bridge"
	networks := ListNetwork()

	for _, network := range networks {
		if network.Name == name {
			return network.ID
		}
	}

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		fmt.Printf("%v\n", err)
		panic(err)
	}
	defer cli.Close()
	ipamConfig := make([]network.IPAMConfig, 0)
	subnet := cniIP + "/16"
	gateway := string([]byte(cniIP)[:len(cniIP)-1]) + "1"
	tmpConfig := network.IPAMConfig{
		Subnet:  subnet,
		Gateway: gateway,
	}
	ipamConfig = append(ipamConfig, tmpConfig)
	ipam := &network.IPAM{
		Config: ipamConfig,
	}
	optionsConfig := make(map[string]string)
	optionsConfig["com.docker.network.bridge.name"] = "miniK8S-bridge"

	networkCreateOption := types.NetworkCreate{
		IPAM:    ipam,
		Options: optionsConfig,
	}
	rep, err := cli.NetworkCreate(context.Background(), "miniK8S-bridge", networkCreateOption)
	fmt.Printf("Netbridge %s has been created\n", rep.ID)

	return rep.ID
}
