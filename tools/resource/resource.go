package resource

import (
	"fmt"
	"github.com/google/cadvisor/client"
	v1 "github.com/google/cadvisor/info/v1"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"mini-kubernetes/tools/def"
	"os/exec"
	"time"
)

func StartCadvisor() (*client.Client, error) {
	arg := fmt.Sprintf("--port=%d", def.CadvisorPort)
	cmd := exec.Command("./cadvisor", arg)
	err := cmd.Start()
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("http://localhost:%d/", def.CadvisorPort)
	return client.NewClient(url)
}

func GetAllContainersInfo(cAdvisorClient *client.Client) []v1.ContainerInfo {
	infoRequest := v1.DefaultContainerInfoRequest()
	containers, _ := cAdvisorClient.AllDockerContainers(&infoRequest)
	return containers
}

func GetContainerInfoByName(cAdvisorClient *client.Client, name string) *v1.ContainerInfo {
	infoRequest := v1.DefaultContainerInfoRequest()
	info, _ := cAdvisorClient.ContainerInfo(name, &infoRequest)
	return info
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
