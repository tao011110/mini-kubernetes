package scheduler_utils

import (
	clientv3 "go.etcd.io/etcd/client/v3"
	"mini-kubernetes/tools/def"
)

func GetNodeResourceInfo(etcdClient *clientv3.Client, nodeID int) *def.NodeResourceSchedulerCache {
	nodeResourceSchedulerCache := def.NodeResourceSchedulerCache{
		Validate: false,
	}
	nodeResource := GetResourceUsageSequenceByNodeID(etcdClient, nodeID)
	//length := len(nodeResource.Sequence)
	//if length != 0 {
	//	latestRecode := nodeResource.Sequence[length-1]
	if nodeResource.Valid {
		latestRecode := nodeResource
		nodeResourceSchedulerCache = def.NodeResourceSchedulerCache{
			CPULoad:     float64(latestRecode.CPULoad) / 1000,
			CPUNum:      latestRecode.CPUNum,
			MemoryUsage: latestRecode.MemoryUsage,
			MemoryTotal: latestRecode.MemoryUsage,
			Validate:    true,
		}
	}
	return &nodeResourceSchedulerCache
}

func GetInfoOfANode(etcdClient *clientv3.Client, nodeID int) (*def.NodeInfoSchedulerCache, []string) {
	nodeInfo := def.NodeInfoSchedulerCache{
		NodeID:          nodeID,
		PodInstanceList: []def.PodInstanceSchedulerCache{},
	}
	replicaNameList := GetAllPodInstancesOfANode(nodeID, etcdClient)
	// for every replica, generate a PodInstanceSchedulerCache
	for _, replicaName := range replicaNameList {
		podInstance := GetPodInstanceByName(etcdClient, replicaName)
		if podInstance.Status != def.FAILED {
			replicaInfo := def.PodInstanceSchedulerCache{
				InstanceName: replicaName,
				PodName:      podInstance.Pod.Metadata.Name,
			}
			nodeInfo.PodInstanceList = append(nodeInfo.PodInstanceList, replicaInfo)
		}
	}
	return &nodeInfo, replicaNameList
}
