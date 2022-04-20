package kubelet

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mini-kubernets/tools/def"
	"mini-kubernets/tools/httpget"
	"time"
)

func sendHeartbeat() {
	request := def.NodeToMasterHeartBeatRequest{
		NodeID:       node.NodeID,
		PodInstances: node.PodInstances,
	}

	body, _ := json.Marshal(request)
	err := httpget.Post("http://" + node.MasterIpAndPort + "/heartbeat").
		ContentType("application/json").
		Body(bytes.NewReader(body)).
		Execute()
	if err != nil {
		fmt.Println(err)
	}
	node.LastHeartbeatSuccessTime = time.Now().Unix()
}
