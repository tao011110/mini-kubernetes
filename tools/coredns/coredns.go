package main

import (
	"bytes"
	"fmt"
	"os/exec"
)

// Dir coredns 执行文件所在目录
var corednsDir = "/etc/coredns"

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

func main() {
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
