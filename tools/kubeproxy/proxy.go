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

	e.DELETE("/delete/clusterIPServiceRule/:clusterIP", deleteCIPServiceRule)

	e.Logger.Fatal(e.Start(":" + strconv.Itoa(ProxyPort)))
}

func addCIPServiceRule(c echo.Context) error {
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

	rules := make([]iptables.Rule, 0)
	for _, pair := range service.PortsBindings {
		num := len(pair.Endpoints)
		i := 1
		for _, endpoint := range pair.Endpoints {
			probability := strconv.FormatFloat(1/(float64(num-i+1)), 'f', -1, 64)
			rule := iptables.Rule{
				Protocol:        pair.Ports.Protocol,
				DestinationIP:   service.IP,
				DestinationPort: strconv.Itoa(int(pair.Ports.Port)),
				DNAT:            endpoint + ":" + pair.Ports.TargetPort,
				Probability:     probability,
				RobinN:          num - i + 1,
			}
			fmt.Printf("add rule is %v\n", rule)
			rules = append(rules, rule)
			fmt.Println(len(rule.Probability))
			//err = ipt.Append("nat", "PREROUTING", "-p", rule.Protocol,
			//	"-d", rule.DestinationIP, "--dport", rule.DestinationPort, "-m", "statistic",
			//	"--mode", "random", "--probability", rule.Probability,
			//	"-j", "DNAT", "--to", rule.DNAT)
			err = ipt.Append("nat", "PREROUTING", "-p", rule.Protocol,
				"-d", rule.DestinationIP, "--dport", rule.DestinationPort, "-m", "statistic",
				"--mode", "nth", "--every", strconv.Itoa(2), "--packet", "1",
				"-j", "DNAT", "--to", rule.DNAT)
			//err = ipt.Append("nat", "OUTPUT", "-p", rule.Protocol,
			//	"-d", rule.DestinationIP, "--dport", rule.DestinationPort, "-j", "DNAT", "--to", rule.DNAT)
			if err != nil {
				fmt.Printf("Add clusterIP service failed: %v", err)
				panic(err)
			}
		}
	}
	newRules[service.IP] = rules
	fmt.Println(len(newRules))
	fmt.Println(newRules[service.IP])
	fmt.Println(service.IP)

	return c.String(200, "Add clusterIP successfully")
}

func deleteCIPServiceRule(c echo.Context) error {
	clusterIP := c.Param("clusterIP")
	fmt.Println("clusterIP is\n" + clusterIP)
	rules := newRules[clusterIP]
	fmt.Println(len(newRules))
	fmt.Println(rules)

	for _, rule := range rules {
		fmt.Printf("delete rule is %v\n", rule)
		//err := ipt.Delete("nat", "PREROUTING", "-p", rule.Protocol,
		//	"-d", rule.DestinationIP, "--dport", rule.DestinationPort, "-m", "statistic",
		//	"--mode", "random", "--probability", rule.Probability,
		//	"-j", "DNAT", "--to", rule.DNAT)
		err := ipt.Delete("nat", "PREROUTING", "-p", rule.Protocol,
			"-d", rule.DestinationIP, "--dport", rule.DestinationPort, "-m", "statistic",
			"--mode", "nth", "--every", strconv.Itoa(rule.RobinN), "--packet", "0",
			"-j", "DNAT", "--to", rule.DNAT)

		if err != nil {
			fmt.Printf("Add clusterIP service failed: %v", err)
		}
	}

	return c.String(200, "Add clusterIP successfully")
}
