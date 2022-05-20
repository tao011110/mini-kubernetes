package kubectl_test

import (
	"mini-kubernetes/tools/kubectl"
	"testing"
)

func Test(t *testing.T) {
	app := kubectl.Initial()
	// brief introduction of kubectl
	err := kubectl.ParseArgs(app, "kubectl")
	if err != nil {
		t.Error("test kubectl fail")
	}

	// test hello
	err = kubectl.ParseArgs(app, "kubectl hello")
	if err != nil {
		t.Error("test kubectl hello fail")
	}

	// test create
	err = kubectl.ParseArgs(app, "kubectl create -f pod.yaml")
	if err != nil {
		t.Error("test kubectl create fail")
	}
}
