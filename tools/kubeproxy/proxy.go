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
var newRules map[string][]iptables.Rule

// ProxyPort kube-proxy所监听的端口号，不建议再进行修改，否则apiserver中也需要相应修改！
var ProxyPort = 3000

func Proxy() {
	newIpt, err := iptables.New()
	if err != nil {
		panic(fmt.Sprintf("New failed: %v", err))
	}
	ipt = newIpt
	newRules = make(map[string][]iptables.Rule)

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
	isExist, err := ipt.ChainExists("nat", "MINI-KUBE-SERVICES")
	if err != nil {
		panic(err)
	}
	if isExist == false {
		fmt.Printf("Do create chain MINI-KUBE-SERVICES")
		err = ipt.NewChain("nat", "MINI-KUBE-SERVICES")
		if err != nil {
			fmt.Printf("Create chain MINI-KUBE-SERVICES failed: %v", err)
			panic(err)
		}
		err = ipt.AppendUnique("nat", "PREROUTING", "-j", "MINI-KUBE-SERVICES")
		if err != nil {
			fmt.Printf("Append rule to MINI-KUBE-SERVICES failed: %v", err)
			panic(err)
		}
	}

	// Create MINI-KUBE-NODEPORTS chain to handle NodePort Service
	isExist, err = ipt.ChainExists("nat", "MINI-KUBE-NODEPORTS")
	if err != nil {
		panic(err)
	}
	if isExist == false {
		fmt.Printf("Do create chain MINI-KUBE-NODEPORTS")
		err = ipt.NewChain("nat", "MINI-KUBE-NODEPORTS")
		if err != nil {
			fmt.Printf("Create chain MINI-KUBE-NODEPORTS failed: %v", err)
			panic(err)
		}
		err = ipt.AppendUnique("nat", "MINI-KUBE-SERVICES", "-j", "MINI-KUBE-NODEPORTS",
			"-m", "addrtype", "--dst-type", "LOCAL")
		if err != nil {
			fmt.Printf("Append rule to MINI-KUBE-NODEPORTS failed: %v", err)
			panic(err)
		}
	}
}

func createSvcChain(clusterIP string) string {
	chainName := "MINI-KUBE-SVC-" + clusterIP
	err := ipt.NewChain("nat", chainName)
	if err != nil {
		fmt.Printf("Create chain %s failed: %v", chainName, err)
		panic(err)
	}
	err = ipt.Insert("nat", "MINI-KUBE-SERVICES", 1,
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

	chainName := createSvcChain(service.IP)

	rules := make([]iptables.Rule, 0)
	for _, pair := range service.PortsBindings {
		num := len(pair.Endpoints)
		i := 0
		for _, endpoint := range pair.Endpoints {
			rule := iptables.Rule{
				Protocol:        pair.Ports.Protocol,
				DestinationIP:   service.IP,
				DestinationPort: strconv.Itoa(int(pair.Ports.Port)),
				DNAT:            endpoint,
				RobinN:          num - i,
			}
			fmt.Printf("add rule is %v\n", rule)
			rules = append(rules, rule)
			err = ipt.AppendUnique("nat", chainName, "-p", rule.Protocol,
				"--dport", rule.DestinationPort, "-m", "statistic",
				"--mode", "nth", "--every", strconv.Itoa(rule.RobinN), "--packet", "0",
				"-j", "DNAT", "--to", rule.DNAT)
			//err = ipt.AppendUnique("nat", "OUTPUT", "-p", rule.Protocol,
			//	"-d", rule.DestinationIP, "--dport", rule.DestinationPort, "-j", "DNAT", "--to", rule.DNAT)
			//err = ipt.AppendUnique("nat", "PREROUTING", "-p", rule.Protocol,
			//	"-d", rule.DestinationIP, "--dport", rule.DestinationPort, "-m", "statistic",
			//	"--mode", "nth", "--every", strconv.Itoa(rule.RobinN), "--packet", "0",
			//	"-j", "DNAT", "--to", rule.DNAT)
			if err != nil {
				fmt.Printf("Add clusterIP service failed: %v", err)
				panic(err)
			}
		}
	}
	newRules[service.IP] = rules
	fmt.Println(newRules[service.IP])

	return c.String(200, "Add clusterIP successfully")
}

func deleteCIPServiceRule(c echo.Context) error {
	clusterIP := c.Param("clusterIP")
	fmt.Println("clusterIP is\n" + clusterIP)

	chainName := "MINI-KUBE-SVC-" + clusterIP

	err := ipt.Delete("nat", "MINI-KUBE-SERVICES",
		"-j", chainName, "-d", clusterIP)
	if err != nil {
		fmt.Printf("Delete ClusterIP service failed: %v", err)
		panic(err)
	}

	err = ipt.ClearAndDeleteChain("nat", chainName)

	//for _, rule := range rules {
	//	fmt.Printf("delete rule is %v\n", rule)
	//	err := ipt.Delete("nat", "PREROUTING", "-p", rule.Protocol,
	//		"-d", rule.DestinationIP, "--dport", rule.DestinationPort, "-m", "statistic",
	//		"--mode", "nth", "--every", strconv.Itoa(rule.RobinN), "--packet", "0",
	//		"-j", "DNAT", "--to", rule.DNAT)
	//
	//}
	if err != nil {
		fmt.Printf("Delete ClusterIP service failed: %v", err)
		panic(err)
	}

	return c.String(200, "Add clusterIP successfully")
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

	chainName := createSvcChain(service.IP)

	rules := make([]iptables.Rule, 0)
	for _, pair := range service.PortsBindings {
		num := len(pair.Endpoints)
		i := 0
		for _, endpoint := range pair.Endpoints {
			rule := iptables.Rule{
				Protocol:        pair.Ports.Protocol,
				DestinationIP:   service.IP,
				DestinationPort: strconv.Itoa(int(pair.Ports.Port)),
				DNAT:            endpoint,
				RobinN:          num - i,
			}
			fmt.Printf("add rule is %v\n", rule)
			rules = append(rules, rule)

			err = ipt.AppendUnique("nat", "MINI-KUBE-NODEPORTS",
				"-d", rule.DestinationIP, "-j", chainName)
			if err != nil {
				fmt.Printf("Add NodePort service failed: %v", err)
				panic(err)
			}

			err = ipt.AppendUnique("nat", chainName, "-p", rule.Protocol,
				"--dport", rule.DestinationPort, "-m", "statistic",
				"--mode", "nth", "--every", strconv.Itoa(rule.RobinN), "--packet", "0",
				"-j", "DNAT", "--to", rule.DNAT)
			if err != nil {
				fmt.Printf("Add NodePort service failed: %v", err)
				panic(err)
			}
		}
	}
	newRules[service.IP] = rules
	fmt.Println(newRules[service.IP])

	return c.String(200, "Add clusterIP successfully")
}

func deleteNPServiceRule(c echo.Context) error {
	clusterIP := c.Param("clusterIP")
	fmt.Println("clusterIP is\n" + clusterIP)

	chainName := "MINI-KUBE-SVC-" + clusterIP

	err := ipt.Delete("nat", "MINI-KUBE-SERVICES",
		"-j", chainName, "-d", clusterIP)
	if err != nil {
		fmt.Printf("Delete NodePort service failed: %v", err)
		panic(err)
	}

	err = ipt.Delete("nat", "MINI-KUBE-NODEPORTS",
		"-d", clusterIP, "-j", chainName)
	if err != nil {
		fmt.Printf("Delete NodePort service failed: %v", err)
		panic(err)
	}

	err = ipt.ClearAndDeleteChain("nat", chainName)

	if err != nil {
		fmt.Printf("Delete NodePort service failed: %v", err)
		panic(err)
	}

	return c.String(200, "Add clusterIP successfully")
}
