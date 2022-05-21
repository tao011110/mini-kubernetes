package coredns

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mini-kubernetes/tools/etcd"
	"os/exec"
	"strings"

	clientv3 "go.etcd.io/etcd/client/v3"
)

// Dir coredns 执行文件所在目录
var corednsDir = "/etc/coredns"

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

// func test_reverse() {
// 	f := func(c rune) bool {
// 		if c == '.' {
// 			return true
// 		} else {
// 			return false
// 		}
// 	}
// 	s := "yingjiu.notr.tech"
// 	result := strings.FieldsFunc(s, f)
// 	fmt.Printf("result:%q\n", result)
// 	key := "/skydns"
// 	for i := len(result) - 1; i >= 0; i-- {
// 		key = key + "/" + result[i]
// 	}
// 	fmt.Printf("result:%s\n", key)
// }

func StartCoredns() {
	cmd := exec.Command("./coredns", "-conf", "Corefile")
	cmd.Dir = corednsDir
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout //标准输出内容到out中
	cmd.Stderr = &stderr //标准输出内容到err中

	err := cmd.Run()
	outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())
	fmt.Printf("out:\n%serr:\n%s", outStr, errStr)

	if err != nil {
		fmt.Printf("[Info] open coredns failed, err:%v\n", err)
	} else {
		fmt.Printf("[Info] open coredns sucess\n")
	}
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
