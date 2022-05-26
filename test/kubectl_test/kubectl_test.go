package kubectl_test

import (
	"mini-kubernetes/tools/kubectl"
	"testing"
)

func TestCreate(t *testing.T) {
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

	// test create pod
	err = kubectl.ParseArgs(app, "kubectl create -f pod.yaml")
	if err != nil {
		t.Error("test kubectl create fail")
	}

	// test delete pod
	err = kubectl.ParseArgs(app, "kubectl delete pod test-pod")
	if err != nil {
		t.Error("test kubectl create fail")
	}
}

func TestDescribe(t *testing.T) {
	app := kubectl.Initial()

	// create pod-1
	err := kubectl.ParseArgs(app, "kubectl create -f pod-1.yaml")
	if err != nil {
		t.Error("test kubectl create fail")
	}

	// create pod-2
	err = kubectl.ParseArgs(app, "kubectl create -f pod-2.yaml")
	if err != nil {
		t.Error("test kubectl create fail")
	}

	// test describe pod
	err = kubectl.ParseArgs(app, "kubectl describe pod test-pod-1")
	if err != nil {
		t.Error("test kubectl describe pod fail")
	}

	// test get pods
	err = kubectl.ParseArgs(app, "kubectl get pods")
	if err != nil {
		t.Error("test kubectl describe pod fail")
	}
}
