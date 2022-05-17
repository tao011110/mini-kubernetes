package kubelet_routines

import (
	"encoding/json"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"mini-kubernetes/tools/def"
	"mini-kubernetes/tools/etcd"
	net_utils "mini-kubernetes/tools/net-utils"
)

func NodesWatch(node *def.Node) {
	fmt.Printf("NodesWatch changes\n")
	prefix := "/node/"
	watchResult := etcd.WatchWithPrefix(node.EtcdClient, prefix)
	for wc := range watchResult {
		//changes := make([]def.Node, 0)
		change := def.Node{}
		added := make([]def.Node, 0)
		deleted := make([]def.Node, 0)
		for _, w := range wc.Events {
			if w.Type == clientv3.EventTypePut {
				fmt.Printf("w.Type is put\n")
				err := json.Unmarshal(w.Kv.Value, &change)
				if err != nil {
					fmt.Println(err)
					panic(err)
				}
				if change.NodeID != node.NodeID {
					// 避免修改node相关参数时，重复PUT导致多次建立隧道而出错
					flag := true
					for _, tmp := range net_utils.NodesList {
						if tmp.NodeID == change.NodeID {
							flag = false
							break
						}
					}
					if flag == true {
						added = append(added, change)
						net_utils.NodesList = append(net_utils.NodesList, change)
					}
				}
			} else {
				if w.Type == clientv3.EventTypeDelete {
					fmt.Printf("w.Type is delete\n")
					fmt.Printf("w.kv.key is %v\n", w.Kv.Key)
					nodeID := 0
					err := json.Unmarshal(w.Kv.Key[6:], &nodeID)
					if err != nil {
						fmt.Println(err)
						panic(err)
					}
					fmt.Printf("nodeID is %v\n", nodeID)
					nodeList := make([]def.Node, 0)
					for _, tmp := range net_utils.NodesList {
						if tmp.NodeID == nodeID && nodeID != node.NodeID {
							deleted = append(deleted, tmp)
						} else {
							nodeList = append(nodeList, tmp)
						}
					}
					net_utils.NodesList = nodeList
				}
			}
		}
		handleNodesChange(added, deleted)
	}
}

func handleNodesChange(adds []def.Node, deletes []def.Node) {
	for _, add := range adds {
		net_utils.CreateVxLan(add)
	}
	for _, _delete := range deletes {
		net_utils.DeleteVxLan(_delete.NodeIP.String())
	}
}
