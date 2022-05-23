package resource_test

import (
	"fmt"
	"github.com/google/cadvisor/client"
	"github.com/robfig/cron"
	"mini-kubernetes/tools/resource"
	"testing"
)

var client_ *client.Client

func Test(t *testing.T) {
	//client_, err := resource.StartCadvisor()
	client_2, err := client.NewClient("http://127.0.0.1:8080")
	if err != nil {
		fmt.Println(err)
	}
	client_ = client_2
	cron2 := cron.New()
	err = cron2.AddFunc("*/30 * * * * *", printResourceInfo)
	if err != nil {
		fmt.Println(err)
	}

	cron2.Start()
	defer cron2.Stop()
	for {

	}
}

func printResourceInfo() {
	//info, _ := client_.MachineInfo()
	//fmt.Printf("%+v\n", info)
	infos := resource.GetAllContainersInfo(client_)
	for _, info := range infos {
		id := info.Id
		name := info.Name
		cpuInfo := info.Stats[len(info.Stats)-1].Cpu.Usage.Total
		memInfo := info.Stats[len(info.Stats)-1].Memory.Usage
		mem := float64(memInfo) / (1024 * 1024)
		cpuUsage := float64(cpuInfo) / (1000 * 1000 * 1000)
		fmt.Printf("id: %s,\nname: %s,\nTotal memoryUsage: %f,\ncpuUasge: %fs\n\n", id, name, mem, cpuUsage)
		fmt.Println(info.Stats[len(info.Stats)-1].Timestamp)
		fmt.Println(info.Stats[len(info.Stats)-1].Timestamp.Second())
		fmt.Println(info.Stats[len(info.Stats)-1].Timestamp.Hour())
		fmt.Println(info.Stats[len(info.Stats)-1].Timestamp.Day())
		fmt.Println(info.Stats[len(info.Stats)-2].Timestamp)
		fmt.Println(info.Stats[len(info.Stats)-2].Timestamp.Second())
		fmt.Println(info.Stats[len(info.Stats)-2].Timestamp.Hour())
		fmt.Println(info.Stats[len(info.Stats)-2].Timestamp.Day())
	}
}
