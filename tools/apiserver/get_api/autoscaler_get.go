package get_api

import (
	"encoding/json"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"mini-kubernetes/tools/controller/controller_utils"
	"mini-kubernetes/tools/def"
	"mini-kubernetes/tools/etcd"
)

func GetAutoscaler(cli *clientv3.Client, autoscalerName string) (*def.AutoscalerDetail, bool) {
	flag := false
	autoscalerKey := def.GetKeyOfAutoscaler(autoscalerName)
	kv := etcd.Get(cli, autoscalerKey).Kvs
	parsedAutoscaler := def.ParsedHorizontalPodAutoscaler{}
	autoscalerValue := make([]byte, 0)
	autoscalerDetail := def.AutoscalerDetail{}
	fmt.Println("autoscalerKey is :   ", autoscalerKey)
	if len(kv) != 0 {
		autoscalerValue = kv[0].Value
		flag = true
		err := json.Unmarshal(autoscalerValue, &parsedAutoscaler)
		if err != nil {
			fmt.Printf("%v\n", err)
			panic(err)
		}

		instancelist := controller_utils.GetReplicaNameListByPodName(cli, parsedAutoscaler.PodName)
		autoscalerDetail.Name = parsedAutoscaler.Name
		autoscalerDetail.CPUMaxValue = parsedAutoscaler.CPUMaxValue
		autoscalerDetail.CPUMinValue = parsedAutoscaler.CPUMinValue
		autoscalerDetail.MemoryMaxValue = parsedAutoscaler.MemoryMaxValue
		autoscalerDetail.MinReplicas = parsedAutoscaler.MinReplicas
		autoscalerDetail.MaxReplicas = parsedAutoscaler.MaxReplicas
		autoscalerDetail.MinReplicas = parsedAutoscaler.MinReplicas
		autoscalerDetail.CreationTimestamp = parsedAutoscaler.StartTime
		autoscalerDetail.CurrentReplicasNum = len(instancelist)
	}

	return &autoscalerDetail, flag
}
