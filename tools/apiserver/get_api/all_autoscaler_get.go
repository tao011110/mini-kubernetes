package get_api

import (
	"encoding/json"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"mini-kubernetes/tools/controller/controller_utils"
	"mini-kubernetes/tools/def"
	"mini-kubernetes/tools/etcd"
	"time"
)

func GetAllAutoscaler(cli *clientv3.Client) ([]def.AutoscalerBrief, bool) {
	parsedAutoscalerKey := "/autoscaler/"
	kvs := etcd.GetWithPrefix(cli, parsedAutoscalerKey).Kvs
	parsedAutoscalerValue := make([]byte, 0)
	parsedAutoscalerBriefList := make([]def.AutoscalerBrief, 0)
	if len(kvs) != 0 {
		for _, kv := range kvs {
			parsedAutoscaler := def.ParsedHorizontalPodAutoscaler{}
			parsedAutoscalerValue = kv.Value
			err := json.Unmarshal(parsedAutoscalerValue, &parsedAutoscaler)
			if err != nil {
				fmt.Printf("%v\n", err)
				panic(err)
			}
			fmt.Println("parsedAutoscaler.Name is " + parsedAutoscaler.Name)
			instancelist := controller_utils.GetReplicaNameListByPodName(cli, parsedAutoscaler.PodName)
			brief := def.AutoscalerBrief{
				Name:     parsedAutoscaler.Name,
				MinPods:  parsedAutoscaler.MinReplicas,
				MaxPods:  parsedAutoscaler.MaxReplicas,
				Age:      time.Now().Sub(parsedAutoscaler.StartTime),
				Replicas: len(instancelist),
			}
			parsedAutoscalerBriefList = append(parsedAutoscalerBriefList, brief)
		}
	} else {
		return parsedAutoscalerBriefList, false
	}

	return parsedAutoscalerBriefList, true
}
