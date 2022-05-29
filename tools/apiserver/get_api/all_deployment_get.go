package get_api

import (
	"encoding/json"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"mini-kubernetes/tools/controller/controller_utils"
	"mini-kubernetes/tools/def"
	"mini-kubernetes/tools/etcd"
	"mini-kubernetes/tools/util"
	"strconv"
	"time"
)

func GetAllDeployment(cli *clientv3.Client) ([]def.DeploymentBrief, bool) {
	deploymentKey := "/deployment/"
	kvs := etcd.GetWithPrefix(cli, deploymentKey).Kvs
	parsedDeploymentValue := make([]byte, 0)
	deploymentBriefList := make([]def.DeploymentBrief, 0)
	if len(kvs) != 0 {
		for _, kv := range kvs {
			parsedDeployment := def.ParsedDeployment{}
			parsedDeploymentValue = kv.Value
			err := json.Unmarshal(parsedDeploymentValue, &parsedDeployment)
			if err != nil {
				fmt.Printf("%v\n", err)
				panic(err)
			}
			fmt.Println("parsedDeployment.Name is " + parsedDeployment.Name)
			instancelist := controller_utils.GetReplicaNameListByPodName(cli, parsedDeployment.PodName)
			fmt.Println("podName is:  ", parsedDeployment.PodName)
			fmt.Println("instancelist is:  ", instancelist)
			health := 0
			ready := 0
			for _, instanceID := range instancelist {
				fmt.Println("instanceID is :   ", instanceID)
				podInstance := util.GetPodInstance(instanceID, cli)
				fmt.Println("status is :   ", podInstance.Status)
				if podInstance.Status != def.FAILED {
					health++
					if podInstance.Status == def.RUNNING {
						ready++
					}
				} else {
					controller_utils.RemovePodInstance(cli, &podInstance)
				}
			}

			brief := def.DeploymentBrief{
				Name:      parsedDeployment.Name,
				Age:       time.Now().Sub(parsedDeployment.StartTime),
				Ready:     strconv.Itoa(ready) + "/" + strconv.Itoa(parsedDeployment.ReplicasNum),
				UpToDate:  health,
				Available: ready,
			}

			deploymentBriefList = append(deploymentBriefList, brief)
		}
	} else {
		return deploymentBriefList, false
	}

	return deploymentBriefList, true
}
