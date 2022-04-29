package pod

import (
	"context"
	"fmt"
	"github.com/docker/docker/client"
	"log"
	"mini-kubernetes/tools/def"
	"mini-kubernetes/tools/docker"
)

func CreateAndStartPod(podInstance *def.PodInstance) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("main here")
	containers := podInstance.Spec.Containers
	containerIDs := make([]string, 0)

	// Create the NetBridge if necessary
	networkID := docker.CreateNetBridge(podInstance.IP)

	// Create the Pause container
	pauseContainerID := docker.CreatePauseContainer(cli, containers, podInstance.Metadata.Name, networkID)

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

		body, err := cli.ContainerCreate(context.Background(), config, hostConfig, networkingConfig, nil, con.Name)
		if err != nil {
			//if error, stop all containers has been created
			podInstance.Status = def.FAILED
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
		podInstance.ContainerStatus[index].Status = def.RUNNING
		podInstance.ContainerStatus[index].ID = body.ID
	}
	podInstance.Status = def.RUNNING
}
