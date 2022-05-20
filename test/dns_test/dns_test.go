package coredns_test

import (
	"mini-kubernetes/tools/coredns"
	"mini-kubernetes/tools/def"
	"mini-kubernetes/tools/etcd"
	"testing"
)

func TestStart(t *testing.T) {
	coredns.StartCoredns()
}

func TestOthers(t *testing.T) {
	cli, err := etcd.Start(def.EtcdDir, def.EtcdPort)
	if err != nil {
		t.Failed()
	}
	defer cli.Close()

	coredns.AddItem(cli, "www.minik8s.com", "127.0.0.1", 1234)
	//coredns.DeleteItem(cli, "www.leffss.com:80")
}
