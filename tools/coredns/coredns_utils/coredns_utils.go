package coredns_utils

import (
	"encoding/json"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"mini-kubernetes/tools/etcd"
	"strings"
)

type Host struct {
	Host string `json:"host"`
}

type Port struct {
	Port string `json:"port"`
}

type HostAndPort struct {
	Host string `json:"host"`
	Port uint16 `json:"port"`
	Ttl  uint8  `json:"ttl"`
}

func AddItem(cli *clientv3.Client, name string, host string, port uint16) {
	/* 基于etcd插件的动态域名增加
	 * E.g 添加一条新的域名解析  yingjiu.notr.tech -> 192.168.1.2
	 * ./etcdctl put /skydns/tech/notr/yingjiu/ '{"host":"192.168.1.2"}'
	 */
	// turn
	f := func(c rune) bool {
		if c == '.' {
			return true
		} else {
			return false
		}
	}
	result := strings.FieldsFunc(name, f)
	key := "/skydns"
	for i := len(result) - 1; i >= 0; i-- {
		key = key + "/" + result[i]
	}

	// set value
	value, err := json.Marshal(HostAndPort{
		Host: host,
		Port: port,
		Ttl:  10,
	})
	if err != nil {
		fmt.Printf("%v\n", err)
		panic(err)
	}

	etcd.Put(cli, key, string(value))
}

func DeleteItem(cli *clientv3.Client, name string) {
	// 基于etcd插件的动态域名删除
	f := func(c rune) bool {
		if c == '.' {
			return true
		} else {
			return false
		}
	}
	result := strings.FieldsFunc(name, f)
	key := "/skydns"
	for i := len(result) - 1; i >= 0; i-- {
		key = key + "/" + result[i]
	}
	etcd.Delete(cli, key)
}
