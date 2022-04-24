package etcd

import (
	"context"
	"fmt"
	"os/exec"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

func Start() {
	// 启动etcd，cmd.Dir填etcd执行文件所在目录
	cmd := exec.Command("./etcd")
	cmd.Dir = "/home/parallels/Downloads/etcd-v3.5.3-linux-arm64"

	if err := cmd.Start(); err != nil {
		fmt.Printf("open etcd failed, err:%v\n", err)
		return
	}
	fmt.Printf("open etcd success\n")

	// 初始化client
	clientConfig := clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	}

	// 创建客户端，连接etcd
	cli, err := clientv3.New(clientConfig)
	if err != nil {
		fmt.Printf("connect to etcd failed, err:%v\n", err)
		return
	}

	timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err = cli.Status(timeoutCtx, clientConfig.Endpoints[0])
	if err != nil {
		fmt.Printf("error checking etcd status: %v\n", err)
		return
	}

	fmt.Println("connect to etcd success")
	defer cli.Close()

	//put get delete watch示例
	/*
		pod, err := yaml.ReadYamlConfig("./yaml/test-2.yaml")
		if err != nil {
			log.Fatal(err)
		}

		// 将对象编码成json数据
		byts, err := json.Marshal(pod)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("struct to json success")

		// put操作
		etcdapi.Put(cli, "testetcdapi", string(byts))

		// get操作
		resp := etcdapi.Get(cli, "testetcdapi") // 获取指定K的值
		if resp == nil {
			fmt.Printf("resp is nil\n")
			return
		}
		for _, ev := range resp.Kvs { // 显示value
			fmt.Printf("key：%s，value：%s\n", ev.Key, ev.Value)
		}

		// delete操作
		etcdapi.Delete(cli, "testetcdapi")

		// watch 程序会卡在这不动 一直监控着watch监控的key值的变化 通过channel管道进行的监控
		watch := etcdapi.Watch(cli, "testetcdapi")
		for wc := range watch {
			for _, w := range wc.Events {
				fmt.Println(string(w.Kv.Key), string(w.Kv.Value), w.Type.String())
			}
		}
	*/
}
