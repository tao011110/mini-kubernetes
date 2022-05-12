package def

import (
	clientv3 "go.etcd.io/etcd/client/v3"
	"sync"
)

type NodeResourceSchedulerCache struct {
	CPULoad     float64
	CPUNum      int
	MemoryUsage uint64
	MemoryTotal uint64
	Validate    bool
}

type PodInstanceSchedulerCache struct {
	InstanceName string
	PodName      string
}

type NodeInfoSchedulerCache struct {
	PodInstanceList []PodInstanceSchedulerCache
	NodeID          int
}

type Scheduler struct {
	ScheduledPodInstancesName []string
	Nodes                     []*NodeInfoSchedulerCache
	EtcdClient                *clientv3.Client
	Lock                      sync.Mutex
	CannotSchedule            []string
	ShouldStop                bool
}
