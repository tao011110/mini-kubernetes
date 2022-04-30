package def

import (
	"github.com/google/cadvisor/client"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	clientv3 "go.etcd.io/etcd/client/v3"
	"mini-kubernetes/tools/pod"
	"net"
)

type NodeResource struct {
	CPUInfo         []cpu.InfoStat        `json:"cpu_info"`
	TotalCPUPercent float64               `json:"total_cpu_percent"`
	PerCPUPercent   []float64             `json:"per_cpu_percent"`
	MemoryInfo      mem.VirtualMemoryStat `json:"memory_info"`
}

type Node struct {
	PodInstances             []*pod.PodInstance
	NodeID                   int
	NodeIP                   net.IP
	NodeName                 string
	MasterIpAndPort          string
	LocalPort                int
	LastHeartbeatSuccessTime int64
	CniIP                    net.IP
	EtcdClient               *clientv3.Client
	CadvisorClient           *client.Client
	ShouldStop               bool
}
