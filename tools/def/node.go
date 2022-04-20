package def

import (
	"net"
)

type Node struct {
	PodInstances    []PodInstance
	NodeID          int
	NodeIP          net.IP
	NodeName        string
	MasterIpAndPort string
	LocalPort       int
}
