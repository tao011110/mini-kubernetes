package etcd

import (
	"context"
	"fmt"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

func Put(cli *clientv3.Client, key string, value string) {
	if cli == nil {
		fmt.Printf("nil client\n")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	_, err := cli.Put(ctx, key, value)
	cancel()

	if err != nil {
		fmt.Printf("put to etcd failed, err:%v\n", err)
		return
	}
}

func Get(cli *clientv3.Client, key string) *clientv3.GetResponse {
	if cli == nil {
		fmt.Printf("nil client\n")
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	resp, err := cli.Get(ctx, key) // 获取指定K的值
	cancel()

	if err != nil {
		fmt.Printf("get to etcd failed")
		return nil
	}
	return resp
}

func GetWithPrefix(cli *clientv3.Client, prefix string) *clientv3.GetResponse {
	if cli == nil {
		fmt.Printf("nil client\n")
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	resp, err := cli.Get(ctx, prefix, clientv3.WithPrefix()) // 获取指定前缀K的值
	cancel()

	if err != nil {
		fmt.Printf("get to etcd failed")
		return nil
	}
	return resp
}

func Delete(cli *clientv3.Client, key string) {
	if cli == nil {
		fmt.Printf("nil client\n")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	_, err := cli.Delete(ctx, key) // 监控一个key的变化
	cancel()

	if err != nil {
		fmt.Println("delete key failed")
	} else {
		fmt.Println("delete key success")
	}
}

func Watch(cli *clientv3.Client, key string) clientv3.WatchChan {
	if cli == nil {
		fmt.Printf("nil client\n")
		return nil
	}

	watch := cli.Watch(context.Background(), key)
	return watch
}

func WatchWithPrefix(cli *clientv3.Client, prefix string) clientv3.WatchChan {
	if cli == nil {
		fmt.Printf("nil client\n")
		return nil
	}

	watch := cli.Watch(context.Background(), prefix, clientv3.WithPrefix())
	return watch
}

func DeleteWithPrefix(cli *clientv3.Client, prefix string) {
	if cli == nil {
		fmt.Printf("nil client\n")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	_, err := cli.Delete(ctx, prefix, clientv3.WithPrefix())
	cancel()

	if err != nil {
		fmt.Println("delete key failed")
	} else {
		fmt.Println("delete key success")
	}
}
