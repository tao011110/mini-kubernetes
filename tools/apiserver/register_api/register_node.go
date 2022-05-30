package register_api

import (
	"encoding/json"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"mini-kubernetes/tools/apiserver/get_api"
	"mini-kubernetes/tools/def"
	"mini-kubernetes/tools/etcd"
	"net"
	"strconv"
	"time"
)

func RegisterNode(cli *clientv3.Client, request def.RegisterToMasterRequest, IpAndPort string) (int, net.IP) {
	{
		nodeList := get_api.GetAllNode(cli)
		for _, node := range nodeList {
			if node.NodeIP.String() == request.LocalIP.String() {
				return node.NodeID, node.CniIP
			}
		}
	}

	registeredNodeID := GetRegisteredNodeID(cli)

	//将新加入集群的node写入到etcd当中
	newFollower := def.Node{}
	newFollower.NodeID = registeredNodeID
	newFollower.NodeIP = request.LocalIP
	newFollower.NodeName = request.NodeName
	newFollower.MasterIpAndPort = IpAndPort
	newFollower.LocalPort = request.LocalPort
	newFollower.ProxyPort = request.ProxyPort
	newFollower.LastHeartbeatSuccessTime = time.Now().Unix()
	newFollower.CniIP = distributeCniIP(registeredNodeID)
	nodeValue, _ := json.Marshal(newFollower)
	etcd.Put(cli, "/node/"+strconv.Itoa(registeredNodeID), string(nodeValue))

	//更新NodeListName
	nodeIDList := make([]int, 0)
	kvs := etcd.Get(cli, def.NodeListID).Kvs
	if len(kvs) != 0 {
		nodeIDListValue := kvs[0].Value
		err := json.Unmarshal(nodeIDListValue, &nodeIDList)
		if err != nil {
			fmt.Printf("%v\n", err)
			panic(err)
		}
	}
	nodeIDList = append(nodeIDList, newFollower.NodeID)
	fmt.Println("newFollower.NodeID is   ", newFollower.NodeID)
	nodeIDListValue, err := json.Marshal(nodeIDList)
	if err != nil {
		fmt.Printf("%v\n", err)
		panic(err)
	}
	fmt.Println("nodeIDList:   ", nodeIDList)
	etcd.Put(cli, def.NodeListID, string(nodeIDListValue))

	return registeredNodeID, newFollower.CniIP
}

func distributeCniIP(registeredNodeID int) net.IP {
	cniIP := net.IPv4(10, 24, byte(registeredNodeID), 0)

	return cniIP
}

func GetRegisteredNodeID(cli *clientv3.Client) int {
	key := "/Persistence/nodeID/"
	kvs := etcd.Get(cli, key).Kvs
	nodeID := 1
	if len(kvs) != 0 {
		nodeIDValue := kvs[0].Value
		err := json.Unmarshal(nodeIDValue, &nodeID)
		if err != nil {
			fmt.Printf("%v\n", err)
			panic(err)
		}
	}
	newNodeID := nodeID + 1
	newNodeIDValue, err := json.Marshal(newNodeID)
	if err != nil {
		fmt.Printf("%v\n", err)
		panic(err)
	}
	etcd.Put(cli, key, string(newNodeIDValue))

	return nodeID
}
