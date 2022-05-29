package resource_test

import (
	"fmt"
	"github.com/google/cadvisor/client"
	"github.com/robfig/cron"
	"mini-kubernetes/tools/resource"
	"testing"
)

var client_ *client.Client

func Test(t *testing.T) {
	//client_, err := resource.StartCadvisor()
	client_2, err := client.NewClient("http://127.0.0.1:8080")
	if err != nil {
		fmt.Println(err)
	}
	client_ = client_2
	cron2 := cron.New()
	err = cron2.AddFunc("*/30 * * * * *", printResourceInfo)
	if err != nil {
		fmt.Println(err)
	}

	cron2.Start()
	defer cron2.Stop()
	for {

	}
}

func printResourceInfo() {
	//info, _ := client_.MachineInfo()
	//fmt.Printf("%+v\n", info)
	infos := resource.GetAllContainersInfo(client_)[0].Stats
	for index, info := range infos {
		fmt.Println("time: ", info.Timestamp, " index: ", index)
	}
}
