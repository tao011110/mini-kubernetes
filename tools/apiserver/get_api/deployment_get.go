package get_api

import (
	"encoding/json"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"mini-kubernetes/tools/def"
	"mini-kubernetes/tools/etcd"
)

func GetDeployment(cli *clientv3.Client, deploymentName string) (*def.Deployment, bool) {
	flag := false
	deploymentKey := def.GetKeyOfDeployment(deploymentName)
	kv := etcd.Get(cli, deploymentKey).Kvs
	deployment := def.Deployment{}
	deploymentValue := make([]byte, 0)
	if len(kv) != 0 {
		deploymentValue = kv[0].Value
		flag = true
		err := json.Unmarshal(deploymentValue, &deployment)
		if err != nil {
			fmt.Printf("%v\n", err)
			panic(err)
		}
	}

	return &deployment, flag
}
