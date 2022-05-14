package net_utils

import (
	"fmt"
	"mini-kubernetes/tools/def"
	"os/exec"
	"strconv"
)

var brNum = 0

type brBinding struct {
	RemoteIp string
	BrName   string
	Subnet   string
}

var brBindings = make([]brBinding, 0)

var NodesList = make([]def.Node, 0)

func CreateVxLan(remoteIp string, subnet string) {
	brName := "br" + strconv.Itoa(brNum)
	binding := brBinding{
		BrName:   brName,
		Subnet:   subnet,
		RemoteIp: remoteIp,
	}
	brBindings = append(brBindings, binding)
	brNum++

	cmd := exec.Command("ovs-vsctl", "add-br", brName)
	err := cmd.Start()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	cmd = exec.Command("brctl", "addif", "miniK8S-bridge", brName)
	err = cmd.Start()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	cmd = exec.Command("ip", "link", "dev", brName, "up")
	err = cmd.Start()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	cmd = exec.Command("ip", "link", "dev", "miniK8S-bridge", "up")
	err = cmd.Start()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	cmd = exec.Command("ip", "route", "add", subnet+"/24", "dev", "miniK8S-bridge")
	err = cmd.Start()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	cmd = exec.Command("ovs-vsctl", "add-port", brName, "vx1", "--", "set",
		"interface", "vx1", "type=vxlan", "options:remote_ip="+remoteIp)
	err = cmd.Start()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}

func DeleteVxLan(remoteIp string) {
	newBrBindings := make([]brBinding, 0)
	for _, binding := range brBindings {
		if remoteIp == binding.RemoteIp {
			cmd := exec.Command("ovs-vsctl", "del-br", binding.BrName)
			err := cmd.Start()
			if err != nil {
				panic(err)
			}
		} else {
			newBrBindings = append(newBrBindings, binding)
		}
	}
}

func InitOVS() {
	// install openvswitch
	cmd := exec.Command("apt", "install", "openvswitch")
	err := cmd.Start()
	if err != nil {
		panic(err)
	}

	// install bridge-utils
	cmd = exec.Command("apt", "install", "bridge-utils")
	err = cmd.Start()
	if err != nil {
		panic(err)
	}
}
