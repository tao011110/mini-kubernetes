package get_api

import (
	"encoding/json"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"mini-kubernetes/tools/def"
	"mini-kubernetes/tools/etcd"
)

func GetAllDNS(cli *clientv3.Client) ([]def.DNSDetail, bool) {
	dnsPrefix := "/DNS/"
	kvs := etcd.GetWithPrefix(cli, dnsPrefix).Kvs
	dnsValue := make([]byte, 0)
	dnsList := make([]def.DNSDetail, 0)
	flag := false
	if len(kvs) != 0 {
		flag = true
		for _, kv := range kvs {
			dns := def.DNSDetail{}
			dnsValue = kv.Value
			err := json.Unmarshal(dnsValue, &dns)
			if err != nil {
				fmt.Printf("%v\n", err)
				panic(err)
			}
			fmt.Println("dns.Name is " + dns.Name)
			dnsList = append(dnsList, dns)
		}
	}

	return dnsList, flag
}
