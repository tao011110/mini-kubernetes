package get_api

import (
	"encoding/json"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"mini-kubernetes/tools/def"
	"mini-kubernetes/tools/etcd"
)

func GetAllNode(cli *clientv3.Client) []def.Node {
	nodeKey := "/node/"
	kvs := etcd.GetWithPrefix(cli, nodeKey).Kvs
	nodeValue := make([]byte, 0)
	nodeList := make([]def.Node, 0)
	if len(kvs) != 0 {
		for _, kv := range kvs {
			node := def.Node{}
			nodeValue = kv.Value
			err := json.Unmarshal(nodeValue, &node)
			if err != nil {
				fmt.Printf("%v\n", err)
				panic(err)
			}
			fmt.Println("node.NodeName is " + node.NodeName)
			nodeList = append(nodeList, node)
		}
	}

	return nodeList
}

func GetAllNodeInfo(cli *clientv3.Client) []def.NodeInfo {
	nodeKey := "/node/"
	kvs := etcd.GetWithPrefix(cli, nodeKey).Kvs
	nodeValue := make([]byte, 0)
	nodeList := make([]def.NodeInfo, 0)
	if len(kvs) != 0 {
		for _, kv := range kvs {
			node := def.Node{}
			nodeValue = kv.Value
			err := json.Unmarshal(nodeValue, &node)
			if err != nil {
				fmt.Printf("%v\n", err)
				panic(err)
			}
			fmt.Println("node.NodeName is " + node.NodeName)
			nodeList = append(nodeList, def.Node2NodeInfo(node))
		}
	}

	return nodeList
}
