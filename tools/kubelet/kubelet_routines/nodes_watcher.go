package kubelet_routines

import (
	"encoding/json"
	"fmt"
	"mini-kubernetes/tools/def"
	"mini-kubernetes/tools/etcd"
	net_utils "mini-kubernetes/tools/net-utils"
)

func NodesWatch(node *def.Node) {
	prefix := "/node/"
	watchResult := etcd.WatchWithPrefix(node.EtcdClient, prefix)
	for wc := range watchResult {
		changes := make([]def.Node, 0)
		change := def.Node{}
		for _, w := range wc.Events {
			err := json.Unmarshal(w.Kv.Value, change)
			if err != nil {
				fmt.Println(err)
				panic(err)
			}
			changes = append(changes, change)
		}
		handleNodesChange(changes)
	}
}

func handleNodesChange(changes []def.Node) {
	adds, deletes := compareNodesList(changes)
	for _, add := range adds {
		net_utils.CreateVxLan(add.NodeIP.String(), add.CniIP.String())
	}
	for _, delete := range deletes {
		net_utils.DeleteVxLan(delete.NodeIP.String())
	}
}

func compareNodesList(changes []def.Node) (added []def.Node, deleted []def.Node) {
	added = make([]def.Node, 0)
	deleted = make([]def.Node, 0)
	for _, change := range changes {
		isAdded := true
		for _, old := range net_utils.NodesList {
			if old.NodeID == change.NodeID {
				isAdded = false
				break
			}
		}
		if isAdded == true {
			added = append(added, change)
		}
	}
	for _, old := range net_utils.NodesList {
		isDeleted := true
		for _, change := range changes {
			if old.NodeID == change.NodeID {
				isDeleted = false
				break
			}
		}
		if isDeleted == true {
			deleted = append(deleted, old)
		}
	}

	return added, deleted
}
