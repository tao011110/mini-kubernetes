package scheduler_utils

import (
	clientv3 "go.etcd.io/etcd/client/v3"
	"mini-kubernetes/tools/def"
)

func NotWithFilter(nodes []int, notWith string, allNodesInfo []*def.NodeInfoSchedulerCache) []int {
	var afterFilterNodes []int
	if notWith == "" {
		return nodes
	}
	for _, node := range nodes {
		in := false
		for _, instance := range allNodesInfo[node].PodInstanceList {
			if instance.PodName == notWith {
				in = true
				break
			}
		}
		if !in {
			afterFilterNodes = append(afterFilterNodes, node)
		}
	}
	return afterFilterNodes
}

func WithFilter(nodes []int, with string, allNodesInfo []*def.NodeInfoSchedulerCache) []int {
	var afterFilterNodes []int
	if with == "" {
		return nodes
	}
	for _, node := range nodes {
		in := false
		for _, instance := range allNodesInfo[node].PodInstanceList {
			if instance.PodName == with {
				in = true
				break
			}
		}
		if in {
			afterFilterNodes = append(afterFilterNodes, node)
		}
	}
	return afterFilterNodes
}

func ResourceFilter(etcdClient *clientv3.Client, nodes []int, CPU int, memory int, allNodesInfo []*def.NodeInfoSchedulerCache) []int {
	var afterFilterNodes []int
	for _, node := range nodes {
		info := GetNodeResourceInfo(etcdClient, allNodesInfo[node].NodeID)
		if info.Validate == false ||
			(info.CPUNum >= CPU && (info.MemoryTotal-info.MemoryUsage) >= uint64(memory)) {
			afterFilterNodes = append(afterFilterNodes, node)
		}
	}
	return afterFilterNodes
}

func ChooseNode(etcdClient *clientv3.Client, nodes []int, allNodesInfo []*def.NodeInfoSchedulerCache) int {
	var chose int
	maxScore := 0
	for _, node := range nodes {
		info := GetNodeResourceInfo(etcdClient, allNodesInfo[node].NodeID)
		score := int(info.MemoryTotal-info.MemoryUsage) + info.CPUNum*int(1000*(1-info.CPULoad))
		if score > maxScore {
			maxScore = score
			chose = node
		}
	}
	return chose
}

// TODO: CPU的单位问题

func MemoryToByte(memString string) int {
	if memString == `0` || memString == `` {
		return 0
	}
	memByte := 0
	for _, c := range memString {
		if c >= '0' && c <= '9' {
			memByte = memByte*10 + int(c-'0')
		} else if c == 'K' || c == 'k' {
			return memByte * 1024
		} else if c == 'M' || c == 'm' {
			return memByte * 1024 * 1024
		} else if c == 'G' || c == 'g' {
			return memByte * 1024 * 1024 * 1024
		}
	}
	return 0
}

func PodResourceRequest(podInstance *def.PodInstance) (int, int) {
	CPU := 0
	memory := 0
	for _, container := range podInstance.Pod.Spec.Containers {
		//requestCPU := container.Resources.ResourceRequest.CPU
		//if requestCPU != 0 && requestCPU > CPU {
		//	CPU = requestCPU
		//}
		memory += MemoryToByte(container.Resources.ResourceRequest.Memory)
	}
	return CPU, memory
}