package util

import (
	"fmt"
	"net"
	"os"
)

func GetLocalIP() net.IP {
	adds, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println(err)
		fmt.Println("cannot get local ip address, exit")
		os.Exit(0)
	}
	for _, address := range adds {
		if ip, flag_ := address.(*net.IPNet); flag_ && !ip.IP.IsLoopback() {
			if ip.IP.To4() != nil {
				return ip.IP.To4()
			}
		}
	}
	os.Exit(0)
	return nil
}
