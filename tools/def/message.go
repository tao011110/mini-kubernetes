package def

import "net"

type RegisterToMasterResponse struct {
	NodeName string
	NodeID   int
}

type RegisterToMasterRequest struct {
	NodeName  string
	LocalIP   net.IP
	LocalPort int
}
