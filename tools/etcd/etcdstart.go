package etcd

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

// Dir etcd执行文件所在目录
//var Dir = "/home/parallels/Downloads/etcd-v3.5.3-linux-arm64"
var Dir = "/home/etcd-v3.2.13-linux-amd64"

// AddToCluster
// kubelet 向 master 注册时调用
// master 将 kubelet 的etcd加入集群
// nodeIPAndPort 格式示例 127.0.0.1:2380
func AddToCluster(nodeName *string, nodeIPAndPort *string) {
	// 调用 ./etcdctl member add 'name' --peer-urls="" 查看状态
	cmd := exec.Command("./etcdctl", "member", "add", *nodeName, "--peer-urls=https://"+*nodeIPAndPort)
	// cmd.Dir 填etcd执行文件所在目录
	cmd.Dir = Dir
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout //标准输出内容到out中
	cmd.Stderr = &stderr //标准输出内容到err中

	err := cmd.Run()
	outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())
	fmt.Printf("./etcdctl member add %s --peer-urls=https://%s\n", *nodeName, *nodeIPAndPort)
	fmt.Printf("out:\n%serr:\n%s", outStr, errStr)

	if err != nil {
		fmt.Printf("[Info] master add new etcd into cluster failed, err:%v\n", err)
	} else {
		fmt.Printf("[Info] master add new etcd into cluster sucess\n")
	}
}

// Start
//	@dir: ""if the same dir
func Start(dir string, etcdPort uint) (*clientv3.Client, error) {

	// 调用 ./etcdctl member list 查看状态
	cmd := exec.Command("./etcdctl", "member", "list")
	// cmd.Dir 填etcd执行文件所在目录
	cmd.Dir = Dir
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout //标准输出内容到out中
	cmd.Stderr = &stderr //标准输出内容到err中

	err := cmd.Run()
	outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())
	fmt.Printf("./etcdctl member list\n")
	fmt.Printf("out:\n%serr:\n%s", outStr, errStr)

	if err != nil {
		fmt.Printf("check member list failed, err:%v\n", err)
		fmt.Printf("[Info] etcd has not started\n")
	} else {
		fmt.Printf("[Info] etcd has started\n")
		goto started
	}

	// etcd 还未启动，尝试启动
	cmd = exec.Command("./etcd")
	cmd.Dir = Dir

	if err := cmd.Start(); err != nil {
		fmt.Printf("open etcd failed, err:%v\n", err)
		return nil, err
	}
	fmt.Printf("open etcd success\n")

started:
	endpoint := fmt.Sprintf("127.0.0.1:%d", etcdPort)
	// 初始化client
	clientConfig := clientv3.Config{
		Endpoints:   []string{endpoint},
		DialTimeout: 5 * time.Second,
	}

	// 创建客户端，连接etcd
	cli, err := clientv3.New(clientConfig)
	if err != nil {
		fmt.Printf("connect to etcd failed, err:%v\n", err)
		return nil, err
	}

	timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err = cli.Status(timeoutCtx, clientConfig.Endpoints[0])
	if err != nil {
		fmt.Printf("error checking etcd status: %v\n", err)
		return nil, err
	}

	fmt.Println("connect to etcd success")
	return cli, nil
}
