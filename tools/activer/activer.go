package main

import (
	"encoding/json"
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/monaco-io/request"
	"github.com/robfig/cron"
	"github.com/thedevsaddam/gojsonq/v2"
	"math"
	"mini-kubernetes/tools/activer/activer_utils"
	"mini-kubernetes/tools/def"
	"mini-kubernetes/tools/etcd"
	"mini-kubernetes/tools/util"
	"net/http"
	"os"
	"time"
)

var activerMeta = def.ActiverCache{
	FunctionsNameList: []string{},
	ShouldStop:        false,
	AccessRecorder:    map[string]int{},
}

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	etcdClient, err := etcd.Start("", def.EtcdPort)
	activerMeta.EtcdClient = etcdClient
	if err != nil {
		e.Logger.Error("Start etcd error!")
		os.Exit(0)
	}
	Initialize()
	go EtcdFunctionsNameListWatcher()
	go AutoExpanderAndShrinker()
	e.GET("/function/:funcname", ProcessFunctionHttpTrigger)
	e.GET("/state_machine/:state_machine_name", ProcessStateMachineHttpTrigger)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", def.ActiverPort)))
}

func Initialize() {
	activerMeta.FunctionsNameList = activer_utils.GetFunctionNameList(activerMeta.EtcdClient)
	for _, functionName := range activerMeta.FunctionsNameList {
		activerMeta.AccessRecorder[functionName] = 0
	}
}

func EtcdFunctionsNameListWatcher() {
	watch := etcd.Watch(activerMeta.EtcdClient, def.FunctionNameListKey)
	for wc := range watch {
		for _, w := range wc.Events {
			var nameList []string
			_ = json.Unmarshal(w.Kv.Value, &nameList)
			HandleFunctionsNameListChange(nameList)
		}
	}
}

func HandleFunctionsNameListChange(functionNameList []string) {
	added, deleted := util.DifferTwoStringList(activerMeta.FunctionsNameList, functionNameList)
	for _, name := range added {
		activerMeta.AccessRecorder[name] = 0
	}
	for _, name := range deleted {
		delete(activerMeta.AccessRecorder, name)
	}
	activerMeta.FunctionsNameList = functionNameList
}

func ProcessFunctionHttpTrigger(c echo.Context) error {
	funcName := c.Param("funcname")
	params := make(map[string]string)
	for k, v := range map[string][]string(c.QueryParams()) {
		if len(v) != 0 {
			params[k] = v[0]
		}
	}
	bytes_ := make([]byte, def.MaxBodySize)
	read, _ := c.Request().Body.Read(bytes_)
	bytes_ = bytes_[:read]
	var body interface{}
	_ = json.Unmarshal(bytes_, &body)
	fmt.Println("params:  ", params)
	fmt.Println("body:   ", body)
	code, response := TriggerFunction(funcName, params, body)
	bytes_, _ = json.Marshal(response)
	return c.String(code, string(bytes_))
}

func TriggerFunction(funcName string, parames map[string]string, body interface{}) (int, interface{}) {
	FlowCount(funcName)
	function := activer_utils.GetFunctionByName(activerMeta.EtcdClient, funcName)
	fmt.Println("funcName:   ", funcName)
	fmt.Println("function:   ", function)
	podReplicaNameList := activer_utils.GetPodReplicaIDListByPodName(activerMeta.EtcdClient, function.PodName)
	service := activer_utils.GetServiceByName(activerMeta.EtcdClient, function.ServiceName)
	if len(podReplicaNameList) == 0 {
		util.AddNPodInstance(function.PodName, 1)
		time.Sleep(30 * time.Second)
		//activer_utils.StartService(function.ServiceName)
	}
	uri := fmt.Sprintf("http://%s:80", service.ClusterIP)
	//uri = fmt.Sprintf("http://10.24.1.2:80")
	fmt.Println(uri)
	c := request.Client{
		URL:    uri,
		Method: "GET",
		Query:  parames,
		JSON:   body,
	}
	var result interface{}
	fmt.Println(c.Send().String())
	resp := c.Send().Scan(&result)
	_ = json.Unmarshal([]byte(c.Send().String()), &result)
	fmt.Println("resp:   ", resp)
	fmt.Println("resp.Response().StatusCode:   ", resp.Response().StatusCode)
	fmt.Println("result:   ", result)
	if resp.Response().StatusCode != 200 {
		return http.StatusInternalServerError, "{}"
	}
	return http.StatusOK, result
}

func ProcessStateMachineHttpTrigger(c echo.Context) error {
	machineName := c.Param("state_machine_name")
	params := make(map[string]string)
	for k, v := range map[string][]string(c.QueryParams()) {
		if len(v) != 0 {
			params[k] = v[0]
		}
	}
	bytes_ := make([]byte, def.MaxBodySize)
	read, _ := c.Request().Body.Read(bytes_)
	bytes_ = bytes_[:read]

	var body interface{}
	_ = json.Unmarshal(bytes_, &body)
	fmt.Println("params:  ", params)
	fmt.Println("body:   ", body)
	code, response := TriggerStateMachine(machineName, params, body)
	bytes_, _ = json.Marshal(response)
	return c.String(code, string(bytes_))
}

func TriggerStateMachine(stateMachineName string, parames map[string]string, body interface{}) (int, interface{}) {
	stateMachine := activer_utils.GetStateMachineByName(activerMeta.EtcdClient, stateMachineName)
	currentState := stateMachine.States[stateMachine.StartAt]
	currentBody := body
	for {
		type_ := gojsonq.New().FromInterface(currentState).Find("Type")
		fmt.Println("type:   ", type_)
		fmt.Println("TriggerStateMachine  currentState:   ", currentState)
		if type_ == "Task" {
			task := def.Task{}
			if currentState.(map[string]interface{})["Next"] != nil {
				task.Type = currentState.(map[string]interface{})["Type"].(string)
				task.Resource = currentState.(map[string]interface{})["Resource"].(string)
				task.Next = currentState.(map[string]interface{})["Next"].(string)
			} else {
				task.Type = currentState.(map[string]interface{})["Type"].(string)
				task.Resource = currentState.(map[string]interface{})["Resource"].(string)
				task.End = currentState.(map[string]interface{})["End"].(bool)
			}
			fmt.Println("get task:  ", task)

			state, response := TriggerFunction(task.Resource, parames, currentBody)
			if state != http.StatusOK || task.End {
				return state, response
			}
			fmt.Println(response)
			currentBody = response
			currentState = stateMachine.States[task.Next]
			fmt.Println(currentState)
		} else if type_ == "Choice" {
			choice := def.Choice{
				Type: currentState.(map[string]interface{})["Type"].(string),
			}
			interfaceMaps := currentState.(map[string]interface{})["Choices"].([]interface{})
			optionList := make([]def.Options, 0)
			for _, interface_ := range interfaceMaps {
				options := def.Options{
					Variable:     interface_.(map[string]interface{})["Variable"].(string),
					StringEquals: interface_.(map[string]interface{})["StringEquals"].(string),
					Next:         interface_.(map[string]interface{})["Next"].(string),
				}
				optionList = append(optionList, options)
			}
			choice.Choices = optionList
			fmt.Println("get choice:   ", choice)
			find := false
			for _, option := range choice.Choices {
				fmt.Println("option.Variable:  ", option.Variable)
				fmt.Println("currentBody:  ", currentBody)
				part := activer_utils.GetPartOfJsonResponce(option.Variable, currentBody)
				fmt.Println("part:  ", part)
				fmt.Println("option.StringEquals:  ", option.StringEquals)
				if part == option.StringEquals {
					currentState = stateMachine.States[option.Next]
					find = true
					break
				}
			}
			if !find {
				fmt.Println("not find")
				return http.StatusInternalServerError, "{}"
			}
		}
	}
}
func AutoExpanderAndShrinker() {
	cron2 := cron.New()
	err := cron2.AddFunc("*/90 * * * * *", ExpandAndShrink)
	if err != nil {
		fmt.Println(err)
	}
	cron2.Start()
	defer cron2.Stop()
	for {
		if activerMeta.ShouldStop {
			break
		}
	}
}

func ExpandAndShrink() {
	//use RCU to avoid lock
	newRecorder := map[string]int{}
	for _, name := range activerMeta.FunctionsNameList {
		newRecorder[name] = 0
	}
	oldRecorder := activerMeta.AccessRecorder
	activerMeta.AccessRecorder = newRecorder
	for _, name := range activerMeta.FunctionsNameList {
		targetReplicaNum := int(math.Ceil(float64(oldRecorder[name]) / 5))
		activer_utils.AdjustReplicaNum2Target(activerMeta.EtcdClient, name, targetReplicaNum)
	}
}

func FlowCount(funcName string) { activerMeta.AccessRecorder[funcName]++ }
