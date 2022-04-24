package cadvisor_client

import (
	"fmt"
	"github.com/google/cadvisor/client"
	v1 "github.com/google/cadvisor/info/v1"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"mini-kubernets/tools/def"
	"os/exec"
	"time"
)

func StartCadvisor() (*client.Client, error) {
	arg := fmt.Sprintf("--port=%d", def.CADVISOR_PORT)
	cmd := exec.Command("./cadvisor", arg)
	err := cmd.Run()
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("http://localhost:%d/", def.CADVISOR_PORT)
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

func GetNodeResourceInfo(cAdvisorClient *client.Client) def.NodeResource {
	totalCPUPercent, _ := cpu.Percent(2*time.Second, false)
	perCPUPercent, _ := cpu.Percent(2*time.Second, true)
	cpuInfo, _ := cpu.Info()
	vmInFo, _ := mem.VirtualMemory()
	NetI
	return def.NodeResource{
		TotalCPUPercent: totalCPUPercent[0],
		PerCPUPercent:   perCPUPercent,
		CPUInfo:         cpuInfo,
		MemoryInfo:      *vmInFo,
	}

}
