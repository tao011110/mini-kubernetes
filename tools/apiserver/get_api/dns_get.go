package get_api

import (
	"encoding/json"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"mini-kubernetes/tools/def"
	"mini-kubernetes/tools/etcd"
)

func GetDNS(cli *clientv3.Client, dnsName string) (def.DNSDetail, bool) {
	flag := false
	dnsKey := "/DNS/" + dnsName
	kv := etcd.Get(cli, dnsKey).Kvs
	dns := def.DNSDetail{}
	dnsValue := make([]byte, 0)
	if len(kv) != 0 {
		dnsValue = kv[0].Value
		flag = true
		err := json.Unmarshal(dnsValue, &dns)
		if err != nil {
			fmt.Printf("%v\n", err)
			panic(err)
		}
	}

	return dns, flag
}
