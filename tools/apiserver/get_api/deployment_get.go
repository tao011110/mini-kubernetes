package get_api

import (
	"encoding/json"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"mini-kubernetes/tools/controller/controller_utils"
	"mini-kubernetes/tools/def"
	"mini-kubernetes/tools/etcd"
	"mini-kubernetes/tools/util"
)

func GetDeployment(cli *clientv3.Client, deploymentName string) (*def.DeploymentDetail, bool) {
	flag := false
	deploymentKey := def.GetKeyOfDeployment(deploymentName)
	kv := etcd.Get(cli, deploymentKey).Kvs
	parsedDeployment := def.ParsedDeployment{}
	parsedDeploymentValue := make([]byte, 0)
	deploymentDetail := def.DeploymentDetail{}
	if len(kv) != 0 {
		parsedDeploymentValue = kv[0].Value
		flag = true
		err := json.Unmarshal(parsedDeploymentValue, &parsedDeployment)
		if err != nil {
			fmt.Printf("%v\n", err)
			panic(err)
		}
		deploymentDetail.Name = parsedDeployment.Name
		deploymentDetail.CreationTimestamp = parsedDeployment.StartTime

		podKey := def.GetKeyOfPod(parsedDeployment.PodName)
		kv := etcd.Get(cli, podKey).Kvs
		pod := def.Pod{}
		podValue := make([]byte, 0)
		if len(kv) != 0 {
			podValue = kv[0].Value
			flag = true
			err := json.Unmarshal(podValue, &pod)
			if err != nil {
				fmt.Printf("%v\n", err)
				panic(err)
			}
		}
		deploymentDetail.PodTemplate = pod

		instancelist := controller_utils.GetReplicaNameListByPodName(cli, parsedDeployment.PodName)
		health := 0
		ready := 0
		for _, instanceID := range instancelist {
			podInstance := util.GetPodInstance(instanceID, cli)
			if podInstance.Status != def.FAILED {
				health++
				if podInstance.Status == def.RUNNING {
					ready++
				}
			} else {
				//controller_utils.RemovePodInstance(cli, &podInstance)
				util.RemovePodInstanceByID(podInstance.ID)
			}
		}
		deploymentDetail.ReplicasState = def.ReplicasState{
			Desired:     parsedDeployment.ReplicasNum,
			Updated:     len(instancelist),
			Total:       len(instancelist),
			Available:   health,
			Unavailable: len(instancelist) - health,
		}
	}

	return &deploymentDetail, flag
}
