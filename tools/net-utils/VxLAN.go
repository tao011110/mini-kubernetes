package net_utils

import (
	"encoding/json"
	"fmt"
	"mini-kubernetes/tools/def"
	"mini-kubernetes/tools/etcd"
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

func CreateVxLan(node def.Node) {
	remoteIp := node.NodeIP.String()
	subnet := node.CniIP.String()
	brName := "br" + strconv.Itoa(brNum)
	vxName := "vx" + strconv.Itoa(brNum)
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

	cmd = exec.Command("ip", "link", "set", "dev", brName, "up")
	err = cmd.Start()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	cmd = exec.Command("ip", "link", "set", "dev", "miniK8S-bridge", "up")
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

	cmd = exec.Command("ovs-vsctl", "add-port", brName, vxName, "--", "set",
		"interface", vxName, "type=vxlan", "options:remote_ip="+remoteIp)
	err = cmd.Start()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	fmt.Printf("CreateVxLan %s to %s, with vx %s\n", binding.BrName, remoteIp, vxName)
	NodesList = append(NodesList, node)
}

func DeleteVxLan(remoteIp string) {
	newBrBindings := make([]brBinding, 0)
	for _, binding := range brBindings {
		if remoteIp == binding.RemoteIp {
			fmt.Printf("DeleteVxLan %s to %s\n", binding.BrName, remoteIp)
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

func InitVxLAN(node *def.Node) {
	nodeKey := "/node/"
	kvs := etcd.GetWithPrefix(node.EtcdClient, nodeKey).Kvs
	nodeValue := make([]byte, 0)
	nodeList := make([]def.Node, 0)
	if len(kvs) != 0 {
		for _, kv := range kvs {
			tmp := def.Node{}
			nodeValue = kv.Value
			err := json.Unmarshal(nodeValue, &tmp)
			if err != nil {
				fmt.Printf("%v\n", err)
				panic(err)
			}
			if tmp.NodeIP.String() != node.NodeIP.String() {
				CreateVxLan(tmp)
				nodeList = append(nodeList, tmp)
			}
		}
	}
	NodesList = nodeList
}

func InitOVS() {
	// install openvswitch
	cmd := exec.Command("apt", "install", "openvswitch-switch")
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
