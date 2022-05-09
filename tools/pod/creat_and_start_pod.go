package pod

import (
	"context"
	"fmt"
	"github.com/docker/docker/client"
	"log"
	"mini-kubernetes/tools/def"
	"mini-kubernetes/tools/docker"
	"mini-kubernetes/tools/util"
	"strings"
	"time"
)

func CreateAndStartPod(podInstance *def.PodInstance, node *def.Node) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		log.Fatal(err)
	}
	defer func(cli *client.Client) {
		_ = cli.Close()
	}(cli)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("main here")
	containers := podInstance.Spec.Containers
	containerIDs := make([]string, 0)

	// Create the NetBridge if necessary
	networkID := docker.CreateNetBridge(node.CniIP.String())

	// Create the Pause container
	pauseContainerID := docker.CreatePauseContainer(cli, containers, podInstance.ID, networkID)
	pauserDetail, _ := docker.InspectContainer(pauseContainerID)
	podInstance.IP = pauserDetail.NetworkSettings.Networks["miniK8S-bridge"].IPAddress
	fmt.Printf("podInstance.IP is %s\n", podInstance.IP)

	for index, con := range containers {
		config := docker.GenerateConfig(con)

		containerMode := "container:" + pauseContainerID
		hostConfig := docker.GenerateHostConfig(con, containerMode)

		tmpCons := make([]def.Container, 0)
		tmpCons = append(tmpCons, con)
		//exportsPort, _ := generatePort(con)
		//fmt.Println(exportsPort)
		//config.ExposedPorts = exportsPort

		networkingConfig := docker.GenerateNetworkingConfig(networkID)

		docker.ImageEnsure(con.Image)
		prefix := podInstance.ID[1:]
		prefix = strings.Replace(prefix, "/", "-", -1)
		name := prefix + `-` + con.Name
		body, err := cli.ContainerCreate(
			context.Background(),
			config, hostConfig,
			networkingConfig,
			nil,
			name)
		if err != nil {
			//if error, stop all containers has been created
			podInstance.Status = def.FAILED
			util.PersistPodInstance(*podInstance, node.EtcdClient)
			for _, id := range containerIDs {
				docker.StopContainer(id)
				_, _ = docker.RemoveContainer(id)
			}
			log.Fatal(err)
			return
		}
		fmt.Println("created " + body.ID)
		containerIDs = append(containerIDs, body.ID)
		docker.StartContainer(body.ID)
		podInstance.ContainerSpec[index].Status = def.RUNNING
		podInstance.ContainerSpec[index].ID = body.ID
		podInstance.ContainerSpec[index].Name = name
		util.PersistPodInstance(*podInstance, node.EtcdClient)
	}
	podInstance.Status = def.RUNNING
	podInstance.StartTime = time.Now()
	util.PersistPodInstance(*podInstance, node.EtcdClient)
	/* 暂时不使用 */
	//go podInstance.PodDaemon()
}
