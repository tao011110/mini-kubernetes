package delete_api

import (
	"encoding/json"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"mini-kubernetes/tools/def"
	"mini-kubernetes/tools/etcd"
	"strconv"
)

func DeleteNode(cli *clientv3.Client, nodeID int) bool {
	// 在etcd中删除node
	fmt.Println("nodeID is :  ", nodeID)
	nodeKey := "/node/" + strconv.Itoa(nodeID)
	resp := etcd.Get(cli, nodeKey)
	if len(resp.Kvs) == 0 {
		return false
	}
	etcd.Delete(cli, nodeKey)

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
	newNodeIDList := make([]int, 0)
	for _, id := range nodeIDList {
		if id != nodeID {
			newNodeIDList = append(newNodeIDList, id)
		}
	}
	nodeIDListValue, err := json.Marshal(newNodeIDList)
	if err != nil {
		fmt.Printf("%v\n", err)
		panic(err)
	}
	fmt.Println("newNodeIDList:   ", newNodeIDList)
	etcd.Put(cli, def.NodeListID, string(nodeIDListValue))

	return true
}
