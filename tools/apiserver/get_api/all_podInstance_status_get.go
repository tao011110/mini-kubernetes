package get_api

import (
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"mini-kubernetes/tools/def"
	"strconv"
	"time"
)

func GetAllPodInstanceStatus(cli *clientv3.Client) ([]def.PodInstanceBrief, bool) {
	podInstanceList, flag := GetAllPodInstance(cli)
	resultList := make([]def.PodInstanceBrief, 0)
	for _, podInstance := range podInstanceList {
		if podInstance.Status != def.SUCCEEDED {
			brief := def.PodInstanceBrief{
				Name:     podInstance.ID[13:],
				Status:   podInstance.Status,
				Restarts: podInstance.RestartCount,
				NodeID:   podInstance.NodeID,
			}
			fmt.Println("podInstance.NodeID,  ", podInstance.NodeID)
			containers := podInstance.ContainerSpec
			count := 0
			for _, container := range containers {
				fmt.Println("container.Status is:  ", container.Status)
				if container.Status == def.RUNNING {
					count++
				} else {
					fmt.Println(container.Status)
				}
			}
			brief.Ready = strconv.Itoa(count) + "/" + strconv.Itoa(len(containers))
			t := time.Now()
			brief.Age = t.Sub(podInstance.StartTime)
			fmt.Println("brief.Age:   ", brief.Age)
			if brief.Age > 31536000*1e9 {
				fmt.Println("brief.Age is too large   ")
				brief.Age = 0
			}

			resultList = append(resultList, brief)
		}
	}
	fmt.Println("resultList:   ", resultList)

	return resultList, flag
}
