package kubelet_utils

import (
	"encoding/json"
	"fmt"
	"github.com/google/cadvisor/client"
	clientv3 "go.etcd.io/etcd/client/v3"
	"mini-kubernetes/tools/def"
	"mini-kubernetes/tools/etcd"
	"mini-kubernetes/tools/resource"
	"time"
)

func RecordNodeResource(nodeID int, etcdClient *clientv3.Client) {
	key := def.KeyNodeResourceUsage(nodeID)
	nodeResource := resource.GetNodeResourceInfo()
	resourceUsage := def.ResourceUsage{
		CPULoad:     int32(nodeResource.TotalCPUPercent * 1000),
		MemoryUsage: nodeResource.MemoryInfo.Used,
		MemoryTotal: nodeResource.MemoryInfo.Total,
		Time:        time.Now(),
		CPUNum:      len(nodeResource.PerCPUPercent),
	}
	byts, _ := json.Marshal(resourceUsage)
	etcd.Put(etcdClient, key, string(byts))
}

func RecordPodInstanceResource(podInstance def.PodInstance, cadvisorClient *client.Client, etcdClient *clientv3.Client) {
	if podInstance.Status != def.RUNNING {
		return
	}
	memoryUsage := uint64(0)
	cpuLoadAverage := int32(0)
	for _, container := range podInstance.ContainerSpec {
		fmt.Println("container.ID is " + container.ID)
		info := resource.GetContainerInfoByID(cadvisorClient, container.ID)
		fmt.Println(info)
		if len(info.Stats) < 2 {
			continue
		}
		memoryUsage += info.Stats[len(info.Stats)-1].Memory.Usage
		time1 := info.Stats[len(info.Stats)-1].Timestamp
		time2 := info.Stats[len(info.Stats)-2].Timestamp
		timeLenth := time1.Unix() - time2.Unix()
		cpuTime1 := info.Stats[len(info.Stats)-1].Cpu.Usage.Total
		cpuTime2 := info.Stats[len(info.Stats)-2].Cpu.Usage.Total
		cpuusageLenth := cpuTime1 - cpuTime2
		cpuNum := (cpuusageLenth * 1000) / uint64(timeLenth) / 1000000
		fmt.Printf("time length is %ds, cpu usage is%dnano\n, memory is %d, last record %s",
			timeLenth, cpuusageLenth, memoryUsage, info.Stats[len(info.Stats)-1].Timestamp.String())
		cpuLoadAverage += int32(cpuNum)
	}
	key := def.GetKeyOfResourceUsageByPodInstanceID(podInstance.ID)
	resourceUsage := def.ResourceUsage{
		CPULoad:     cpuLoadAverage,
		MemoryUsage: memoryUsage,
		Time:        time.Now(),
	}
	byts, _ := json.Marshal(resourceUsage)
	etcd.Put(etcdClient, key, string(byts))
}
