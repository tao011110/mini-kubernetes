package resource_test

import (
	"fmt"
	"mini-kubernets/tools/resource"
	"testing"
)

func Test(t *testing.T) {
	nodeResource := resource.GetNodeResourceInfo()
	fmt.Printf("%+v\n", nodeResource)
	client_, err := resource.StartCadvisor()
	if err != nil {
		fmt.Println(err)
	}
	info, _ := client_.MachineInfo()
	fmt.Printf("%+v\n", info)
}
