package delete_api

import (
	"encoding/json"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"mini-kubernetes/tools/def"
	"mini-kubernetes/tools/etcd"
)

func DeleteDeployment(cli *clientv3.Client, deploymentName string) bool {
	// 在deployment_list_name 中删除 deployment's name
	{
		deploymentListNameKey := def.DeploymentListName
		deploymentListNameValue := etcd.Get(cli, deploymentListNameKey).Kvs[0].Value
		list := make([]string, 0)
		err := json.Unmarshal(deploymentListNameValue, &list)
		if err != nil {
			fmt.Printf("%v\n", err)
			panic(err)
		}
		newList := make([]string, 0)
		for _, name := range list {
			if name != deploymentName {
				newList = append(newList, name)
			}
		}
		deploymentListNameValue, err = json.Marshal(newList)
		if err != nil {
			fmt.Printf("%v\n", err)
			panic(err)
		}
		etcd.Put(cli, deploymentListNameKey, string(deploymentListNameValue))
	}

	// 在etcd中删除ParsedDeployment
	{
		deploymentKey := def.GetKeyOfDeployment(deploymentName)
		resp := etcd.Get(cli, deploymentKey)
		if len(resp.Kvs) == 0 {
			return false
		}
		etcd.Delete(cli, deploymentKey)
	}

	// 在etcd中删除ParsedDeployment的pod和podInstance
	{
		podPrefix := "/pod/" + deploymentName
		etcd.DeleteWithPrefix(cli, podPrefix)
		//podInstancePrefix := "/podInstance/" + deploymentName
		//etcd.DeleteWithPrefix(cli, podInstancePrefix)
	}

	return true
}
