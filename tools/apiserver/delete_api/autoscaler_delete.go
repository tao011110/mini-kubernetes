package delete_api

import (
	"encoding/json"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"mini-kubernetes/tools/def"
	"mini-kubernetes/tools/etcd"
)

func DeleteAutoscaler(cli *clientv3.Client, autoscalerName string) bool {
	// 在autoscaler_list_name 中删除 autoscaler's name
	{
		autoscalerListNameKey := def.HorizontalPodAutoscalerListName
		autoscalerListNameValue := etcd.Get(cli, autoscalerListNameKey).Kvs[0].Value
		list := make([]string, 0)
		err := json.Unmarshal(autoscalerListNameValue, &list)
		if err != nil {
			fmt.Printf("%v\n", err)
			panic(err)
		}
		newList := make([]string, 0)
		for _, name := range list {
			if name != autoscalerName {
				newList = append(newList, name)
			}
		}
		autoscalerListNameValue, err = json.Marshal(newList)
		if err != nil {
			fmt.Printf("%v\n", err)
			panic(err)
		}
		etcd.Put(cli, autoscalerListNameKey, string(autoscalerListNameValue))
	}

	// 在etcd中删除ParsedAutoscaler
	{
		autoscalerKey := def.GetKeyOfAutoscaler(autoscalerName)
		resp := etcd.Get(cli, autoscalerKey)
		if len(resp.Kvs) == 0 {
			return false
		}
		etcd.Delete(cli, autoscalerKey)
	}

	// 在etcd中删除autoscaler的pod
	{
		podPrefix := "/pod/" + autoscalerName
		etcd.DeleteWithPrefix(cli, podPrefix)
		//podInstancePrefix := "/podInstance/" + autoscalerName
		//etcd.DeleteWithPrefix(cli, podInstancePrefix)
	}

	return true
}
