package kubectl_test

import (
	"mini-kubernetes/tools/kubectl"
	"testing"
	"time"
)

func TestPod(t *testing.T) {
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

func TestServiceAndDNS(t *testing.T) {
	app := kubectl.Initial()

	// create service-1
	err := kubectl.ParseArgs(app, "kubectl create -f service-1.yaml")
	if err != nil {
		t.Error("test kubectl create fail")
	}

	// create service-2
	err = kubectl.ParseArgs(app, "kubectl create -f service-2.yaml")
	if err != nil {
		t.Error("test kubectl create fail")
	}

	// test describe service-1
	err = kubectl.ParseArgs(app, "kubectl describe service service-1")
	if err != nil {
		t.Error("test kubectl describe service fail")
	}

	// test get services
	err = kubectl.ParseArgs(app, "kubectl get services")
	if err != nil {
		t.Error("test kubectl get services fail")
	}

	// create dns
	err = kubectl.ParseArgs(app, "kubectl create -f dns.yaml")
	if err != nil {
		t.Error("test kubectl create dns fail")
	}

	// test describe dns
	err = kubectl.ParseArgs(app, "kubectl describe dns test-dns")
	if err != nil {
		t.Error("test kubectl describe dns fail")
	}

	// test get dns
	err = kubectl.ParseArgs(app, "kubectl get dns")
	if err != nil {
		t.Error("test kubectl get dns fail")
	}

	// test delete service-1
	err = kubectl.ParseArgs(app, "kubectl delete service service-1")
	if err != nil {
		t.Error("test kubectl delete service fail")
	}

	// test delete service-2
	err = kubectl.ParseArgs(app, "kubectl delete service service-2")
	if err != nil {
		t.Error("test kubectl delete service fail")
	}
}

func TestDeployment(t *testing.T) {
	app := kubectl.Initial()
	
	// create deployment
	err := kubectl.ParseArgs(app, "kubectl create -f ../yaml_test/dep.yaml")
	if err != nil {
		t.Error("test kubectl create fail")
	}

	// test describe deployment
	err = kubectl.ParseArgs(app, "kubectl describe deployment deployment-hjk")
	if err != nil {
		t.Error("test kubectl describe deployment fail")
	}

	// test get deployment
	err = kubectl.ParseArgs(app, "kubectl get deployment")
	if err != nil {
		t.Error("test kubectl get deployment fail")
	}

	// test delete deployment
	err = kubectl.ParseArgs(app, "kubectl delete deployment deployment-hjk")
	if err != nil {
		t.Error("test kubectl delete deployment fail")
	}
}

func TestAutoscaler(t *testing.T) {
	app := kubectl.Initial()
	
	// create autoscaler
	err := kubectl.ParseArgs(app, "kubectl create -f ../yaml_test/auto.yaml")
	if err != nil {
		t.Error("test kubectl create fail")
	}

	// test describe autoscaler
	err = kubectl.ParseArgs(app, "kubectl describe autoscaler mynginx")
	if err != nil {
		t.Error("test kubectl describe autoscaler fail")
	}

	// test get autoscaler
	err = kubectl.ParseArgs(app, "kubectl get autoscaler")
	if err != nil {
		t.Error("test kubectl get autoscaler fail")
	}

	// test delete autoscaler
	err = kubectl.ParseArgs(app, "kubectl delete autoscaler mynginx")
	if err != nil {
		t.Error("test kubectl delete autoscaler fail")
	}
}

func TestGPU(t *testing.T) {
	app := kubectl.Initial()
	
	// create gpujob
	err := kubectl.ParseArgs(app, "kubectl create -f ../apiserver_test/gpu_test.yaml")
	if err != nil {
		t.Error("test kubectl create fail")
	}

	time.Sleep(50 * time.Second)

	// test describe gpujob
	err = kubectl.ParseArgs(app, "kubectl describe gpujob GPUJob-test")
	if err != nil {
		t.Error("test kubectl describe gpujob fail")
	}

	// test get gpujob
	err = kubectl.ParseArgs(app, "kubectl get gpujob")
	if err != nil {
		t.Error("test kubectl get gpujob fail")
	}
}

func TestFuncAndState(t *testing.T) {
	app := kubectl.Initial()
	
	// create function
	err := kubectl.ParseArgs(app, "kubectl create -f ../apiserver_test/state1.yaml")
	err = kubectl.ParseArgs(app, "kubectl create -f ../apiserver_test/state2.yaml")
	err = kubectl.ParseArgs(app, "kubectl create -f ../apiserver_test/state3.yaml")
	err = kubectl.ParseArgs(app, "kubectl create -f ../apiserver_test/state4.yaml")
	if err != nil {
		t.Error("test kubectl create function fail")
	}

	// test describe function
	err = kubectl.ParseArgs(app, "kubectl describe function state1_function")
	if err != nil {
		t.Error("test kubectl describe function fail")
	}

	// test get function
	err = kubectl.ParseArgs(app, "kubectl get function")
	if err != nil {
		t.Error("test kubectl get function fail")
	}

	// create statemachine
	err = kubectl.ParseArgs(app, "kubectl create -f ../yaml_test/workflow.json")
	if err != nil {
		t.Error("test kubectl create statemachine fail")
	}

	// test describe statemachine
	err = kubectl.ParseArgs(app, "kubectl describe statemachine test_state_machine")
	if err != nil {
		t.Error("test kubectl describe statemachine fail")
	}

	// test get statemachine
	err = kubectl.ParseArgs(app, "kubectl get statemachine")
	if err != nil {
		t.Error("test kubectl get statemachine fail")
	}
	
	// test delete statemachine
	err = kubectl.ParseArgs(app, "kubectl delete statemachine test_state_machine")
	if err != nil {
		t.Error("test kubectl delete statemachine fail")
	}

	// test delete functions
	err = kubectl.ParseArgs(app, "kubectl delete function state1_function")
	err = kubectl.ParseArgs(app, "kubectl delete function state2_function")
	err = kubectl.ParseArgs(app, "kubectl delete function state3_function")
	err = kubectl.ParseArgs(app, "kubectl delete function state4_function")
	if err != nil {
		t.Error("test kubectl delete function fail")
	}
}

