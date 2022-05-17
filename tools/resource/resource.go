package resource

import (
	"fmt"
	"github.com/google/cadvisor/client"
	v1 "github.com/google/cadvisor/info/v1"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"mini-kubernetes/tools/def"
	"time"
)

func StartCadvisor() (*client.Client, error) {
	// arg := fmt.Sprintf("--port=%d", def.CadvisorPort)
	// cmd := exec.Command("./cadvisor", arg)
	// err := cmd.Start()
	// if err != nil {
	// 	return nil, err
	// }
	url := fmt.Sprintf("http://localhost:%d/", def.CadvisorPort)
	return client.NewClient(url)
}

func GetAllContainersInfo(cAdvisorClient *client.Client) []v1.ContainerInfo {
	infoRequest := v1.DefaultContainerInfoRequest()
	containers, err := cAdvisorClient.AllDockerContainers(&infoRequest)
	if err != nil {
		print(err)
	}
	return containers
}

func GetContainerInfoByID(cAdvisorClient *client.Client, id string) *v1.ContainerInfo {
	infos := GetAllContainersInfo(cAdvisorClient)
	for _, info := range infos {
		if info.Id == id {
			return &info
		}
	}
	return nil
}

func GetNodeResourceInfo() def.NodeResource {
	totalCPUPercent, _ := cpu.Percent(2*time.Second, false)
	perCPUPercent, _ := cpu.Percent(2*time.Second, true)
	cpuInfo, _ := cpu.Info()
	vmInFo, _ := mem.VirtualMemory()
	return def.NodeResource{
		TotalCPUPercent: totalCPUPercent[0],
		PerCPUPercent:   perCPUPercent,
		CPUInfo:         cpuInfo,
		MemoryInfo:      *vmInFo,
	}
}
