package docker

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"log"
	"mini-kubernets/tools/def"
	"mini-kubernets/tools/yaml"
	"strconv"
)

// Create the Pause container, which acts as the parent of all containers in the pod
func createPauseContainer(cli *client.Client, cons []def.Container, podName string) string {
	ImageEnsure("registry.aliyuncs.com/google_containers/pause")
	config := &container.Config{
		Image: "registry.aliyuncs.com/google_containers/pause",
	}
	hostConfig := &container.HostConfig{}
	exportsPort, portMap := generatePorts(cons)
	config.ExposedPorts = exportsPort
	hostConfig.PortBindings = portMap

	body, err := cli.ContainerCreate(context.Background(), config, hostConfig, nil, nil, "pause-"+podName)
	if err != nil {
		log.Fatal(err)
	}
	StartContainer(body.ID)

	return body.ID
}

func CreateContrainer(path string) []string {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	pod, err := yaml.ReadYamlConfig(path)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("main here")
	containers := pod.Spec.Containers
	containerIDs := make([]string, 0)

	// Create the Pause container
	pauseContainerID := createPauseContainer(cli, containers, pod.Metadata.Name)

	for _, con := range containers {
		config := generateConfig(con)

		containerMode := "container:" + pauseContainerID
		hostConfig := generateHostConfig(con, containerMode)

		tmpCons := make([]def.Container, 0)
		tmpCons = append(tmpCons, con)
		//exportsPort, _ := generatePort(con)
		//fmt.Println(exportsPort)
		//config.ExposedPorts = exportsPort

		ImageEnsure(con.Image)

		body, err := cli.ContainerCreate(context.Background(), config, hostConfig, nil, nil, con.Name)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("created " + body.ID)
		containerIDs = append(containerIDs, body.ID)
	}

	return containerIDs
}

// generate Config
func generateConfig(con def.Container) *container.Config {
	config := &container.Config{
		Image:      con.Image,
		WorkingDir: con.WorkingDir,
	}

	if len(con.Commands) != 0 {
		config.Entrypoint = con.Commands
	}
	if len(con.Args) != 0 {
		config.Cmd = con.Args
	}

	return config
}

// generate HostConfig
func generateHostConfig(con def.Container, containerMode string) *container.HostConfig {
	resourcesConfig := container.Resources{}

	limits := con.Resources.ResourceLimit
	// get the CPU limits
	cpuLimits := 0
	cpus := 0
	if len(limits.CPU) != 0 {
		if limits.CPU[len(limits.CPU)-1] == 'm' {
			// 'm' means millicore, eg: '500m' means 0.5 logical CPU
			cpuLimits, _ = strconv.Atoi(string([]byte(limits.CPU)[:len(limits.CPU)-1]))
			// Convert the number of configured cpus to nanocPU
			cpus = cpuLimits * 1e6
		} else {
			cpuLimits, _ = strconv.Atoi(string([]byte(limits.CPU)[:len(limits.CPU)]))
			// Convert the number of configured cpus to nanocPU
			cpus = cpuLimits * 1e9
		}
		resourcesConfig.NanoCPUs = int64(cpus)
	}

	// get the Memory limits
	if len(limits.Memory) != 0 {
		memoryLimits := 0
		memory := 0
		if limits.Memory[len(limits.Memory)-1] == 'i' {
			// Ki, Mi, Gi --1024
			switch limits.Memory[len(limits.Memory)-2] {
			case 'K':
				memoryLimits, _ = strconv.Atoi(string([]byte(limits.Memory)[:len(limits.Memory)-2]))
				memory = memoryLimits << 10
			case 'M':
				memoryLimits, _ = strconv.Atoi(string([]byte(limits.Memory)[:len(limits.Memory)-2]))
				memory = memoryLimits << 20
			case 'G':
				memoryLimits, _ = strconv.Atoi(string([]byte(limits.Memory)[:len(limits.Memory)-2]))
				memory = memoryLimits << 30
			}
		} else {
			// K，M，G -- 1000
			switch limits.Memory[len(limits.Memory)-1] {
			case 'K':
				memoryLimits, _ = strconv.Atoi(string([]byte(limits.Memory)[:len(limits.Memory)-1]))
				memory = memoryLimits * 10e3
			case 'M':
				memoryLimits, _ = strconv.Atoi(string([]byte(limits.Memory)[:len(limits.Memory)-1]))
				memory = memoryLimits * 10e6
			case 'G':
				memoryLimits, _ = strconv.Atoi(string([]byte(limits.Memory)[:len(limits.Memory)-1]))
				memory = memoryLimits * 10e9
			}
		}
		resourcesConfig.Memory = int64(memory)
	}

	// get the mounted volumes
	mounts := make([]mount.Mount, 0)
	volumes := con.VolumeMounts
	for _, volume := range volumes {
		mountVolume := mount.Mount{
			Type:   "volume",
			Source: volume.Name,
			Target: volume.MountPath,
		}
		mounts = append(mounts, mountVolume)
	}

	hostConfig := &container.HostConfig{
		Resources:   resourcesConfig,
		Mounts:      mounts,
		PidMode:     container.PidMode(containerMode),
		IpcMode:     container.IpcMode(containerMode),
		NetworkMode: container.NetworkMode(containerMode),
	}

	return hostConfig
}

// get exposedPorts and hostPorts
func generatePorts(cons []def.Container) (nat.PortSet, nat.PortMap) {
	exportPorts := make(nat.PortSet)
	portMap := make(nat.PortMap)
	for _, con := range cons {
		portsMappings := con.PortMappings
		for _, ports := range portsMappings {
			if ports.ContainerPort != 0 {
				port, err := nat.NewPort(ports.Protocol, strconv.Itoa(int(ports.ContainerPort)))
				if err != nil {
					log.Fatal(err)
				}
				exportPorts[port] = struct{}{}

				if ports.HostPort != 0 {
					portBind := nat.PortBinding{HostPort: strconv.Itoa(int(ports.HostPort))}
					tmp := make([]nat.PortBinding, 0, 1)
					tmp = append(tmp, portBind)
					portMap[port] = tmp
				}
			}
		}
	}

	return exportPorts, portMap
}
