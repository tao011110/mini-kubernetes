package etcd_test

import (
	"encoding/json"
	"fmt"
	"log"
	"mini-kubernetes/tools/def"
	"mini-kubernetes/tools/etcd"
	"mini-kubernetes/tools/yaml"
	"testing"
)

// Testing `go test` to run
func Test(t *testing.T) {
	cli, err := etcd.Start(etcd.Dir, def.EtcdPort)
	if err != nil {
		t.Failed()
	}
	defer cli.Close()

	//put get delete watch示例
	pod, err := yaml.ReadPodYamlConfig("./yaml/test-2.yaml")
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
	etcd.Put(cli, "testetcdapi", string(byts))

	// get操作
	resp := etcd.Get(cli, "testetcdapi") // 获取指定K的值

	if resp == nil {
		fmt.Printf("resp is nil\n")
		return
	}
	for _, ev := range resp.Kvs { // 显示value
		fmt.Printf("key：%s，value：%s\n", ev.Key, ev.Value)
	}

	// delete操作
	etcd.Delete(cli, "testetcdapi")

	// watch 程序会卡在这不动 一直监控着watch监控的key值的变化 通过channel管道进行的监控
	watch := etcd.Watch(cli, "testetcdapi")
	for wc := range watch {
		for _, w := range wc.Events {
			fmt.Println(string(w.Kv.Key), string(w.Kv.Value), w.Type.String())
		}
	}

}
