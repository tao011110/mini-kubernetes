package register_api

import (
	"encoding/json"
	clientv3 "go.etcd.io/etcd/client/v3"
	"mini-kubernetes/tools/def"
	"mini-kubernetes/tools/etcd"
	"net"
	"strconv"
	"time"
)

var registeredNodeID = 0

func RegisterNode(cli *clientv3.Client, request def.RegisterToMasterRequest, IpAndPort string) (int, net.IP) {
	registeredNodeID++

	//将新加入集群的node写入到etcd当中
	newFollower := def.Node{}
	newFollower.NodeID = registeredNodeID
	newFollower.NodeIP = request.LocalIP
	newFollower.NodeName = request.NodeName
	newFollower.MasterIpAndPort = IpAndPort
	newFollower.LocalPort = request.LocalPort
	newFollower.ProxyPort = request.ProxyPort
	newFollower.LastHeartbeatSuccessTime = time.Now().Unix()
	newFollower.CniIP = distributeCniIP()
	nodeValue, _ := json.Marshal(newFollower)
	etcd.Put(cli, "/node/"+strconv.Itoa(registeredNodeID), string(nodeValue))

	//将 kubelet 的etcd加入etcd集群中
	nodeIPAndPort := newFollower.NodeIP.String() + ":2380"

	etcd.AddToCluster(&newFollower.NodeName, &nodeIPAndPort)

	return registeredNodeID, newFollower.CniIP
}

func distributeCniIP() net.IP {
	cniIP := net.IPv4(10, 24, byte(registeredNodeID), 0)

	return cniIP
}
