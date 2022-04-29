package delete_api

import (
	clientv3 "go.etcd.io/etcd/client/v3"
	"mini-kubernetes/tools/etcd"
)

func DeletePod(cli *clientv3.Client, podName string) {
	//在etcd中删除podInstance
	podInstanceKey := "/podInstance/" + podName
	etcd.Delete(cli, podInstanceKey)

	//在etcd中删除pod
	//TODO: 暂定删除podInstance时也一并删除pod，将来的实现可以会进行修改
	podKey := "/pod/" + podName
	etcd.Delete(cli, podKey)
}
