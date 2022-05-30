package function_api

import (
	"encoding/json"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"mini-kubernetes/tools/def"
	"mini-kubernetes/tools/etcd"
	"mini-kubernetes/tools/util"
)

func DeleteFunction(cli *clientv3.Client, functionName string) (bool, string) {
	idPrefix := fmt.Sprintf("functional_%s", functionName)
	podName := fmt.Sprintf("pod_%s", idPrefix)
	serviceName := fmt.Sprintf("service_%s", idPrefix)

	// 从function_name_list中删除该function
	{
		var functionList []string
		util.EtcdUnmarshal(etcd.Get(cli, def.FunctionNameListKey), &functionList)
		newFunctionList := make([]string, 0)
		for _, function := range functionList {
			if function != functionName {
				newFunctionList = append(newFunctionList, function)
			}
		}
		value, err := json.Marshal(newFunctionList)
		if err != nil {
			fmt.Printf("%v\n", err)
			panic(err)
		}
		etcd.Put(cli, def.FunctionNameListKey, string(value))
	}

	// 删除function
	{
		functionKey := def.GetKeyOfFunction(functionName)
		resp := etcd.Get(cli, functionKey)
		if len(resp.Kvs) == 0 {
			//return false, ""
		}
		etcd.Delete(cli, functionKey)
	}

	{ // 删除之前创建的functionPodInstance
		replicasNameListPrefix := def.GetKeyOfPodReplicasNameListByPodName(podName)
		var instanceIDList []string
		resp := etcd.GetWithPrefix(cli, replicasNameListPrefix)
		fmt.Println("resp.Kvs:  ", resp.Kvs)
		if len(resp.Kvs) != 0 {
			value := resp.Kvs[0].Value
			fmt.Println("value:  ", value)
			err := json.Unmarshal(value, &instanceIDList)
			if err != nil {
				fmt.Printf("%v\n", err)
				panic(err)
			}
			fmt.Println("replicasNameListPrefix:   ", replicasNameListPrefix)
			fmt.Println("instanceIDList is:   ", instanceIDList)
			for _, instanceID := range instanceIDList {
				fmt.Println("When delete function, try to delete podInstance:  ", instanceID)
				FuncDeletePodInstance(cli, instanceID)
			}

			// 删除对应的replicas_name_list
			etcd.DeleteWithPrefix(cli, replicasNameListPrefix)
		}
	}

	// 删除对应的pod
	{
		etcd.DeleteWithPrefix(cli, def.GetKeyOfPod(podName))
	}

	// 删除对应的service
	servicePrefix := def.GetKeyOfService(serviceName)
	clusterIP := ""
	resp := etcd.GetWithPrefix(cli, servicePrefix)
	etcd.Delete(cli, servicePrefix)

	serviceValue := resp.Kvs[0].Value
	service := def.Service{}
	err := json.Unmarshal(serviceValue, &service)
	clusterIP = service.ClusterIP
	if err != nil {
		fmt.Printf("%v\n", err)
		panic(err)
	}

	return true, clusterIP
}
