package get_api

import (
	"encoding/json"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"mini-kubernetes/tools/apiserver/apiserver_utils"
	"mini-kubernetes/tools/def"
	"mini-kubernetes/tools/etcd"
)

func GetPodInstance(cli *clientv3.Client, podInstanceName string) (*def.PodInstance, bool) {
	flag := false
	podInstanceKey := "/podInstance/" + podInstanceName
	kv := etcd.Get(cli, podInstanceKey).Kvs
	podInstance := def.PodInstance{}
	podInstanceValue := make([]byte, 0)
	if len(kv) != 0 {
		podInstanceValue = kv[0].Value
		flag = true
		err := json.Unmarshal(podInstanceValue, &podInstance)
		if err != nil {
			fmt.Printf("%v\n", err)
			panic(err)
		}
	}

	// add for heartbeat
	if apiserver_utils.GetNodeByID(cli, podInstance.NodeID).Status == def.NotReady {
		podInstance.Status = def.UNKNOWN
	}

	return &podInstance, flag
}
