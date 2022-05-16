package resource_test

import (
	"fmt"
	"github.com/google/cadvisor/client"
	"mini-kubernetes/tools/resource"
	"testing"
)

func Test(t *testing.T) {
	//client_, err := resource.StartCadvisor()
	client_, err := client.NewClient("http://127.0.0.1:8080")
	if err != nil {
		fmt.Println(err)
	}
	//info, _ := client_.MachineInfo()
	//fmt.Printf("%+v\n", info)
	info := resource.GetAllContainersInfo(client_)[0]
	id := info.Id
	name := info.Name
	cpuInfo := info.Stats[len(info.Stats)-1].Cpu.Usage.Total
	memInfo := info.Stats[len(info.Stats)-1].Memory.Usage
	mem := float64(memInfo) / (1024 * 1024)
	cpuLoad := float64(cpuInfo) / (1000 * 1000 * 1000)
	fmt.Printf("id: %s,\nname: %s,\nTotal memoryUsage: %f,\ncpuUasge: %fs\n", id, name, mem, cpuLoad)
}
