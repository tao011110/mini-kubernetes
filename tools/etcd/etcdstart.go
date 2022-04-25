package etcd

import (
	"context"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"os/exec"
	"time"
)

// Start
//	@dir: ""if the same dir
func Start(dir string, etcdPort uint) (*clientv3.Client, error) {
	// 启动etcd，cmd.Dir填etcd执行文件所在目录
	cmd := exec.Command("./etcd")
	if dir != "" {
		cmd.Dir = dir
	}

	if err := cmd.Start(); err != nil {
		fmt.Printf("open etcd failed, err:%v\n", err)
		return nil, err
	}
	fmt.Printf("open etcd success\n")

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
