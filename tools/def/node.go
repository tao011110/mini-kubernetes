package def

import (
	"github.com/google/cadvisor/client"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	clientv3 "go.etcd.io/etcd/client/v3"
	"net"
)

type NodeResource struct {
	CPUInfo         []cpu.InfoStat        `json:"cpu_info"`
	TotalCPUPercent float64               `json:"total_cpu_percent"`
	PerCPUPercent   []float64             `json:"per_cpu_percent"`
	MemoryInfo      mem.VirtualMemoryStat `json:"memory_info"`
}

const (
	Ready    = 0
	NotReady = 1
)

type Node struct {
	PodInstances             []*PodInstance
	NodeID                   int `json:"node_id"`
	NodeIP                   net.IP
	NodeName                 string
	MasterIpAndPort          string
	LocalPort                int
	ProxyPort                int
	LastHeartbeatSuccessTime int64
	CniIP                    net.IP
	EtcdClient               *clientv3.Client
	CadvisorClient           *client.Client
	ShouldStop               bool
	Status                   int
}

type NodeInfo struct {
	NodeID                   int    `json:"node_id"`
	NodeIP                   net.IP `json:"node_ip"`
	NodeName                 string `json:"node_name"`
	MasterIpAndPort          string `json:"master_ip_and_port"`
	LocalPort                int    `json:"local_port"`
	ProxyPort                int    `json:"proxy_port"`
	LastHeartbeatSuccessTime int64  `json:"last_heartbeat_success_time"`
	CniIP                    net.IP `json:"cni_ip"`
	Status                   int    `json:"status"`
}

func Node2NodeInfo(node Node) NodeInfo {
	return NodeInfo{
		NodeID:                   node.NodeID,
		NodeIP:                   node.NodeIP,
		NodeName:                 node.NodeName,
		MasterIpAndPort:          node.MasterIpAndPort,
		LocalPort:                node.LocalPort,
		ProxyPort:                node.ProxyPort,
		LastHeartbeatSuccessTime: node.LastHeartbeatSuccessTime,
		CniIP:                    node.CniIP,
		Status:                   node.Status,
	}
}
