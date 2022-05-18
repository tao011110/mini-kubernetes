package get_api

import (
	"encoding/json"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"mini-kubernetes/tools/def"
	"mini-kubernetes/tools/etcd"
)

func GetAllDeployment(cli *clientv3.Client) ([]def.Deployment, bool) {
	deploymentKey := "/deployment/"
	kvs := etcd.GetWithPrefix(cli, deploymentKey).Kvs
	deploymentValue := make([]byte, 0)
	deploymentList := make([]def.Deployment, 0)
	if len(kvs) != 0 {
		for _, kv := range kvs {
			deployment := def.Deployment{}
			deploymentValue = kv.Value
			err := json.Unmarshal(deploymentValue, &deployment)
			if err != nil {
				fmt.Printf("%v\n", err)
				panic(err)
			}
			fmt.Println("deployment.Metadata.Name is " + deployment.Metadata.Name)
			deploymentList = append(deploymentList, deployment)
		}
	} else {
		return deploymentList, false
	}

	return deploymentList, true
}
