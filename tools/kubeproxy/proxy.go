package kubeproxy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"mini-kubernetes/tools/def"
	"mini-kubernetes/tools/iptables"
	"strconv"
)

var ipt *iptables.IPTables
var svc2sep map[string][]string

// ProxyPort kube-proxy所监听的端口号，不建议再进行修改，否则apiserver中也需要相应修改！
var ProxyPort = 3000

func Proxy() {
	newIpt, err := iptables.New()
	if err != nil {
		panic(fmt.Sprintf("New failed: %v", err))
	}
	ipt = newIpt

	svc2sep = make(map[string][]string)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/add/clusterIPServiceRule", addCIPServiceRule)
	e.POST("/add/nodePortServiceRule", addNPServiceRule)

	e.DELETE("/delete/clusterIPServiceRule/:clusterIP", deleteCIPServiceRule)
	e.DELETE("/delete/nodePortServiceRule/:clusterIP", deleteNPServiceRule)

	e.Logger.Fatal(e.Start(":" + strconv.Itoa(ProxyPort)))
}

func initTable() {
	isExist, err := ipt.ChainExists("nat", "mK8S-SERVICES")
	if err != nil {
		panic(err)
	}
	if isExist == false {
		fmt.Printf("Do create chain mK8S-SERVICES")
		err = ipt.NewChain("nat", "mK8S-SERVICES")
		if err != nil {
			fmt.Printf("Create chain mK8S-SERVICES failed: %v", err)
			panic(err)
		}
		err = ipt.AppendUnique("nat", "PREROUTING", "-j", "mK8S-SERVICES")
		if err != nil {
			fmt.Printf("Append rule to mK8S-SERVICES failed: %v", err)
			panic(err)
		}
	}

	// Create mK8S-NODEPORTS chain to handle NodePort Service
	isExist, err = ipt.ChainExists("nat", "mK8S-NODEPORTS")
	if err != nil {
		panic(err)
	}
	if isExist == false {
		fmt.Printf("Do create chain mK8S-NODEPORTS")
		err = ipt.NewChain("nat", "mK8S-NODEPORTS")
		if err != nil {
			fmt.Printf("Create chain mK8S-NODEPORTS failed: %v", err)
			panic(err)
		}
		err = ipt.AppendUnique("nat", "mK8S-SERVICES", "-j", "mK8S-NODEPORTS",
			"-m", "addrtype", "--dst-type", "LOCAL")
		if err != nil {
			fmt.Printf("Append rule to mK8S-NODEPORTS failed: %v", err)
			panic(err)
		}
	}
}

func createSvcChain(clusterIP string) string {
	chainName := "mK8S-SVC-" + clusterIP
	err := ipt.NewChain("nat", chainName)
	if err != nil {
		fmt.Printf("Create chain %s failed: %v", chainName, err)
		panic(err)
	}

	// Add svc chain into mK8S-SERVICES chain, according to its clusterIP
	err = ipt.Insert("nat", "mK8S-SERVICES", 1,
		"-j", chainName, "-d", clusterIP)
	if err != nil {
		fmt.Printf("Append rule to %s failed: %v", chainName, err)
		panic(err)
	}

	return chainName
}

func addCIPServiceRule(c echo.Context) error {
	initTable()

	service := &def.Service{}
	requestBody := new(bytes.Buffer)
	_, err := requestBody.ReadFrom(c.Request().Body)
	if err != nil {
		fmt.Printf("%v\n", err)
		panic(err)
	}

	err = json.Unmarshal(requestBody.Bytes(), &service)
	if err != nil {
		fmt.Printf("%v\n", err)
		panic(err)
	}

	// Create svc chain
	svcName := createSvcChain(service.IP)

	sepList := make([]string, 0)

	for _, pair := range service.PortsBindings {
		protocol := pair.Ports.Protocol
		destinationIP := service.IP
		destinationPort := strconv.Itoa(int(pair.Ports.Port))

		// Create sep chain
		sepName := "mK8S-SEP-" + destinationIP + "-" + destinationPort
		err := ipt.NewChain("nat", sepName)
		if err != nil {
			fmt.Printf("Create chain %s failed: %v", sepName, err)
			panic(err)
		}
		sepList = append(sepList, sepName)

		// Add sep chain into svc chain, according to its destinationPort
		err = ipt.AppendUnique("nat", svcName, "-p", protocol,
			"--dport", destinationPort, "-j", sepName)
		if err != nil {
			fmt.Printf("Add SEP chain %s failed: %v\n", sepName, err)
			panic(err)
		}

		num := len(pair.Endpoints)
		i := 0
		fmt.Println(pair.Endpoints)

		for _, endpoint := range pair.Endpoints {
			// Fill in the sep chain with endpoints
			err = ipt.AppendUnique("nat", sepName, "-p", protocol,
				"-m", "statistic", "--mode", "nth", "--every", strconv.Itoa(num-i),
				"--packet", "0", "-j", "DNAT", "--to", endpoint)
			if err != nil {
				fmt.Printf("Add NodePort service failed: %v", err)
				panic(err)
			}
			i++
		}
	}
	svc2sep[svcName] = sepList

	return c.String(200, "Add clusterIP successfully")
}

func deleteCIPServiceRule(c echo.Context) error {
	clusterIP := c.Param("clusterIP")
	fmt.Println("clusterIP is\n" + clusterIP)

	svcName := "mK8S-SVC-" + clusterIP
	// Clear and delete all sep chains
	sepList := svc2sep[svcName]
	for _, sepName := range sepList {
		err := ipt.ClearAndDeleteChain("nat", sepName)
		if err != nil {
			fmt.Printf("Delete ClusterIP service failed: %v", err)
			panic(err)
		}
	}

	// Delete svc rule in mK8S-SERVICES chain
	err := ipt.Delete("nat", "mK8S-SERVICES",
		"-j", svcName, "-d", clusterIP)
	if err != nil {
		fmt.Printf("Delete ClusterIP service failed: %v", err)
		panic(err)
	}

	// Clear and delete svc chain
	err = ipt.ClearAndDeleteChain("nat", svcName)
	if err != nil {
		fmt.Printf("Delete ClusterIP service failed: %v", err)
		panic(err)
	}

	return c.String(200, "Delete clusterIP successfully")
}

func addNPServiceRule(c echo.Context) error {
	fmt.Println("addNPServiceRule")
	initTable()

	service := &def.Service{}
	requestBody := new(bytes.Buffer)
	_, err := requestBody.ReadFrom(c.Request().Body)
	if err != nil {
		fmt.Printf("%v\n", err)
		panic(err)
	}

	err = json.Unmarshal(requestBody.Bytes(), &service)
	if err != nil {
		fmt.Printf("%v\n", err)
		panic(err)
	}

	// Create svc chain
	svcName := createSvcChain(service.IP)

	sepList := make([]string, 0)

	for _, pair := range service.PortsBindings {
		protocol := pair.Ports.Protocol
		destinationIP := service.IP
		destinationPort := strconv.Itoa(int(pair.Ports.Port))
		nodePort := strconv.Itoa(int(pair.Ports.NodePort))

		// Create sep chain
		sepName := "mK8S-SEP-" + destinationIP + "-" + destinationPort
		err := ipt.NewChain("nat", sepName)
		if err != nil {
			fmt.Printf("Create chain %s failed: %v", sepName, err)
			panic(err)
		}
		sepList = append(sepList, sepName)

		// Add sep chain into svc chain, according to its destinationPort
		err = ipt.AppendUnique("nat", svcName, "-p", protocol,
			"--dport", destinationPort, "-j", sepName)
		if err != nil {
			fmt.Printf("Add SEP chain %s failed: %v\n", sepName, err)
			panic(err)
		}

		// Add svc chain into mK8S-NODEPORTS chain, according to its nodePort
		err = ipt.AppendUnique("nat", "mK8S-NODEPORTS", "-p", protocol,
			"--dport", nodePort, "-j", svcName)
		if err != nil {
			fmt.Printf("Add NodePort service failed: %v", err)
			panic(err)
		}

		num := len(pair.Endpoints)
		i := 0

		for _, endpoint := range pair.Endpoints {
			// Fill in the sep chain with endpoints
			err = ipt.AppendUnique("nat", sepName, "-p", protocol,
				"-m", "statistic", "--mode", "nth", "--every", strconv.Itoa(num-i),
				"--packet", "0", "-j", "DNAT", "--to", endpoint)
			if err != nil {
				fmt.Printf("Add NodePort service failed: %v", err)
				panic(err)
			}
			i++
		}
	}
	svc2sep[svcName] = sepList

	return c.String(200, "Add clusterIP successfully")
}

func deleteNPServiceRule(c echo.Context) error {
	clusterIP := c.Param("clusterIP")
	fmt.Println("clusterIP is\n" + clusterIP)

	svcName := "mK8S-SVC-" + clusterIP
	// Clear and delete all sep chains
	sepList := svc2sep[svcName]
	for _, sepName := range sepList {
		err := ipt.ClearAndDeleteChain("nat", sepName)
		if err != nil {
			fmt.Printf("Delete NodePort service failed: %v", err)
			panic(err)
		}
	}

	// Delete svc rule in mK8S-SERVICES chain
	err := ipt.Delete("nat", "mK8S-SERVICES",
		"-j", svcName, "-d", clusterIP)
	if err != nil {
		fmt.Printf("Delete NodePort service failed: %v", err)
		panic(err)
	}

	// Delete svc rule in mK8S-NODEPORTS chain
	err = ipt.Delete("nat", "mK8S-NODEPORTS",
		"-d", clusterIP, "-j", svcName)
	if err != nil {
		fmt.Printf("Delete NodePort service failed: %v", err)
		panic(err)
	}

	// Clear and delete svc chain
	err = ipt.ClearAndDeleteChain("nat", svcName)
	if err != nil {
		fmt.Printf("Delete NodePort service failed: %v", err)
		panic(err)
	}

	return c.String(200, "Add clusterIP successfully")
}
