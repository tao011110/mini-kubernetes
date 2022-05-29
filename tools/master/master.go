package main

import (
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"mini-kubernetes/tools/apiserver"
	"mini-kubernetes/tools/def"
	"mini-kubernetes/tools/etcd"
	"mini-kubernetes/tools/util"
)

func main() {
	//创建etcd client
	cli, err := etcd.Start(def.EtcdDir, def.EtcdPort)
	if err != nil {
		fmt.Printf("%v\n", err)
		panic(err)
	}
	defer func(cli *clientv3.Client) {
		err := cli.Close()
		if err != nil {

		}
	}(cli)

	//启动apiserver, 注意Linux / Unix系统默认规定,低端口号(1-1024),user组是不能访问的,需要root组才行
	apiserver.Start(util.GetLocalIP().String(), fmt.Sprintf("%d", def.MasterPort), cli)
}
