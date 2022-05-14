package master

import (
	"fmt"
	"mini-kubernetes/tools/apiserver"
	"mini-kubernetes/tools/def"
	"mini-kubernetes/tools/etcd"
)

var IP = "192.168.1.7"
var Port = "8000"

func Start() {
	//创建etcd client
	cli, err := etcd.Start(etcd.Dir, def.EtcdPort)
	if err != nil {
		fmt.Printf("%v\n", err)
		panic(err)
	}
	defer cli.Close()

	//启动apiserver, 注意Linux / Unix系统默认规定,低端口号(1-1024),user组是不能访问的,需要root组才行
	apiserver.Start(IP, Port, cli)
}
