package net_utils_test

import (
	net_utils "mini-kubernetes/tools/net-utils"
	"testing"
)

func Test(t *testing.T) {
	net_utils.InitOVS()
	net_utils.CreateVxLan("10.119.11.111", "10.24.3.0")
	//net_utils.DeleteVxLan("10.119.11.111")
}
