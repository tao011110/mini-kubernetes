package def

import (
	clientv3 "go.etcd.io/etcd/client/v3"
	"sync"
)

type ControllerMeta struct {
	EtcdClient         *clientv3.Client
	ParsedDeployments  []*ParsedDeployment
	DeploymentNameList []string
	Lock               sync.Mutex
	ShouldStop         bool
}
