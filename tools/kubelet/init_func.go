package kubelet

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"mini-kubernetes/tools/def"
	"mini-kubernetes/tools/httpget"
	"net"
	"os"
)

/*
	command format:./kubelet --name `nodeName` --master `masterIP:port` --port `localPort`
	for example: ./kubelet --name node1 --master 192.168.55.184:80 --port 80
*/
func parseArgs(nodeName *string, masterIPAndPort *string, localPort *int) {
	flag.StringVar(nodeName, "--name", "undefined", "name of the node, `node+nodeIP` by default")
	flag.StringVar(masterIPAndPort, "--master", "undefined", "name of the node, `node+nodeIP` by default")
	flag.IntVar(localPort, "--port", 80, "local port to communicate with master")
	flag.Parse()
	/*
		TODO: Check IP format legality
	*/
	if *masterIPAndPort == "undefined" {
		fmt.Println("Master Ip And Port Error!")
		os.Exit(0)
	}
}

/*
	get local Ip
*/
func getLocalIP() net.IP {
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

/*
	register node to master, using http post
*/
func registerToMaster(node *def.Node) error {
	response := def.RegisterToMasterResponse{}
	request := def.RegisterToMasterRequest{
		NodeName:  node.NodeName,
		LocalIP:   node.NodeIP,
		LocalPort: node.LocalPort,
		ProxyPort: node.ProxyPort,
	}

	body, _ := json.Marshal(request)
	err, _ := httpget.Post("http://" + node.MasterIpAndPort + "/register_node").
		ContentType("application/json").
		Body(bytes.NewReader(body)).
		GetJson(&response).
		Execute()
	if err != nil {
		fmt.Println(err)
		return err
	}
	node.NodeID = response.NodeID
	node.NodeName = response.NodeName
	node.CniIP = response.CniIP
	return nil
}
