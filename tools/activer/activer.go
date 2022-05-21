package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/robfig/cron"
	"github.com/thedevsaddam/gojsonq/v2"
	"mini-kubernetes/tools/activer/activer_utils"
	"mini-kubernetes/tools/def"
	"mini-kubernetes/tools/etcd"
	"mini-kubernetes/tools/httpget"
	"mini-kubernetes/tools/util"
	"net/http"
	"os"
)

//NOTE: 根据pod对应的replica数目来判断集群中是否有实例, replica数目降为0时service删除, 冷启动时service创建

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
	go AutoExpanderAndShrinker()
	go EtcdFunctionsNameListWatcher()
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

func ProcessStateMachineHttpTrigger(c echo.Context) error {
	machineName := c.Param("state_machine_name")
	parames := c.QueryParams().Encode()
	bytes_ := make([]byte, def.MaxBodySize)
	read, _ := c.Request().Body.Read(bytes_)
	bytes_ = bytes_[:read]
	body := string(bytes_)
	return c.String(TriggerStateMachine(machineName, parames, body))
}

func ProcessFunctionHttpTrigger(c echo.Context) error {
	funcName := c.Param("funcname")
	parames := c.QueryParams().Encode()
	bytes_ := make([]byte, def.MaxBodySize)
	read, _ := c.Request().Body.Read(bytes_)
	bytes_ = bytes_[:read]
	body := string(bytes_)
	return c.String(TriggerFunction(funcName, parames, body))
}

func TriggerStateMachine(stateMachineName string, parames string, body string) (int, string) {
	stateMachine := activer_utils.GetStateMachineByName(activerMeta.EtcdClient, stateMachineName)
	currentState := stateMachine.States[stateMachine.StartAt]
	currentBody := body
	for {
		type_ := gojsonq.New().FromString(currentState).Find("Type")
		if type_ == "Task" {
			task := def.Task{}
			_ = json.Unmarshal([]byte(currentState), &task)
			state, responce := TriggerFunction(task.Resource, parames, currentBody)
			if state != http.StatusOK || task.End {
				return state, responce
			}
			currentBody = responce
			currentState = stateMachine.States[task.Next]
		} else if type_ == "choice" {
			choice := def.Choice{}
			_ = json.Unmarshal([]byte(currentState), &choice)
			find := false
			for _, option := range choice.Choices {
				part := activer_utils.GetPartOfJsonResponce(option.Variable, currentBody)
				if part == option.StringEquals {
					currentState = stateMachine.States[option.Next]
					find = true
					break
				}
			}
			if !find {
				return http.StatusInternalServerError, ""
			}
		}

	}
}

func TriggerFunction(funcName string, parames string, body string) (int, string) {
	FlowCount(funcName)
	function := activer_utils.GetFunctionByName(activerMeta.EtcdClient, funcName)
	podReplicaNameList := activer_utils.GetPodReplicaIDListByPodName(activerMeta.EtcdClient, function.PodName)
	service := activer_utils.GetServiceByName(activerMeta.EtcdClient, function.ServiceName)
	if len(podReplicaNameList) == 0 {
		activer_utils.AddNPodInstance(function.PodName, 1)
		activer_utils.StartService(function.ServiceName)
	}
	uri := fmt.Sprintf("%s:80?%s", service.ClusterIP, parames)
	response := ""
	err, status := httpget.Post(uri).
		ContentType("application/json").
		Body(bytes.NewReader([]byte(body))).
		GetString(&response).
		Execute()
	if err != nil || status != "200 OK" {
		return http.StatusInternalServerError, ""
	}
	return http.StatusOK, response
}

func AutoExpanderAndShrinker() {
	cron2 := cron.New()
	err := cron2.AddFunc("*/30 * * * * *", ExpandAndShrink)
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
		targetReplicaNum := oldRecorder[name] / 100
		activer_utils.AdjustReplicaNum2Target(activerMeta.EtcdClient, name, targetReplicaNum)
	}
}

func FlowCount(funcName string) { activerMeta.AccessRecorder[funcName]++ }
