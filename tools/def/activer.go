package def

import clientv3 "go.etcd.io/etcd/client/v3"

type ActiverCache struct {
	FunctionsNameList []string
	EtcdClient        *clientv3.Client
	ShouldStop        bool
	AccessRecorder    map[string]int
}
