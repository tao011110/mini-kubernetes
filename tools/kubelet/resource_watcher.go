package kubelet

import (
	"encoding/json"
	"fmt"
	"github.com/robfig/cron"
	"mini-kubernetes/tools/def"
	"mini-kubernetes/tools/etcd"
	"mini-kubernetes/tools/resource"
	"time"
)

func ResourceMonitoring() {
	cron2 := cron.New()
	err := cron2.AddFunc("*/30 * * * * *", recordResource)
	if err != nil {
		fmt.Println(err)
	}

	cron2.Start()
	defer cron2.Stop()
	for {
		if node.ShouldStop {
			break
		}
	}
}

func recordResource() {
	for _, podInstance := range node.PodInstances {
		if podInstance.Status != def.RUNNING {
			continue
		}
		memoryUsage := uint64(0)
		cpuLoadAverage := int32(0)
		for _, container := range podInstance.ContainerSpec {
			info := resource.GetContainerInfoByName(node.CadvisorClient, container.Name)
			memoryUsage += info.Stats[len(info.Stats)-1].Memory.Usage
			cpuLoadAverage += info.Stats[len(info.Stats)-1].Cpu.LoadAverage
		}
		key := fmt.Sprintf("%s_resource_usage", podInstance.ID)
		resp := etcd.Get(node.EtcdClient, key)
		resourceSeq := def.ResourceUsageSequence{}
		jsonString := ``
		for _, ev := range resp.Kvs {
			jsonString += fmt.Sprintf(`"%s":"%s", `, ev.Key, ev.Value)
		}
		jsonString = fmt.Sprintf(`{%s}`, jsonString)
		_ = json.Unmarshal([]byte(jsonString), &resourceSeq)
		resourceUsage := def.ResourceUsage{
			CPULoad:     cpuLoadAverage,
			MemoryUsage: memoryUsage,
			Time:        time.Now(),
		}
		if len(resourceSeq.Sequence) < 30 {
			resourceSeq.Sequence = append(resourceSeq.Sequence, resourceUsage)
		} else {
			resourceSeq.Sequence = append(resourceSeq.Sequence[1:], resourceUsage)
		}
		byts, _ := json.Marshal(resourceSeq)
		etcd.Put(node.EtcdClient, key, string(byts))
	}
	key := fmt.Sprintf("%d_resource_usage", node.NodeID)
	resp := etcd.Get(node.EtcdClient, key)
	resourceSeq := def.ResourceUsageSequence{}
	jsonString := ``
	for _, ev := range resp.Kvs {
		jsonString += fmt.Sprintf(`"%s":"%s"`, ev.Key, ev.Value)
	}
	jsonString = fmt.Sprintf(`{%s}`, jsonString)
	_ = json.Unmarshal([]byte(jsonString), &resourceSeq)
	nodeResource := resource.GetNodeResourceInfo()
	resourceUsage := def.ResourceUsage{
		CPULoad:     int32(nodeResource.TotalCPUPercent * 1000),
		MemoryUsage: nodeResource.MemoryInfo.Used,
		MemoryTotal: nodeResource.MemoryInfo.Total,
		Time:        time.Now(),
	}
	if len(resourceSeq.Sequence) < 30 {
		resourceSeq.Sequence = append(resourceSeq.Sequence, resourceUsage)
	} else {
		resourceSeq.Sequence = append(resourceSeq.Sequence[1:], resourceUsage)
	}
	byts, _ := json.Marshal(resourceSeq)
	etcd.Put(node.EtcdClient, key, string(byts))
}
