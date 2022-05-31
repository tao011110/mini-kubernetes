GO_CMD=go
GO_BUILD=$(GO_CMD) build
GO_CLEAN=$(GO_CMD) clean
GO_TEST=$(GO_CMD) test
GO_GET=$(GO_CMD) get
TARGET_ACTIVER=activer
TARGET_KUBELET=kubelet
TARGET_APISERVER=master
TARGET_CONTROLLER=controller
TARGET_KUBECTL=kubectl
TARGET_SCHEDULER=scheduler
TARGET_PROXY=proxy
TARGET_CADVISOR=cadvisor
.DEFAULT_GOAL := default

GO_TEST_PATH= './test/yaml_test'
#			'./test/application_yaml_config_test' \
#			'./test/django_test' \
#			'./test/dns_test' \
#			'./test/docker_test' \
#			'./test/echo_test' \
#			'./test/etcd_test' \
#			'./test/image_factory_test' \
#			'./test/kubectl_test' \
#			'./test/kubeproxy_test' \
#			'./test/master_test' \
#			'./test/net_utils_test' \
#			'./test/resource_test' \
#			'./test/slurmGenrator_test' \
			''./test/apiserver_test''
.PHONY:test

all: test master node

master: kubectl apiServer controller scheduler activer

node: kubelet proxy

default: master node

test:
	for path in ${GO_TEST_PATH}; do \
	$(GO_TEST) $$path -v ; \
	done

kubectl:
	$(GO_BUILD) -o ./bin/$(TARGET_KUBECTL) ./tools/kubectl/kubectl.go

apiServer:
	$(GO_BUILD) -o ./bin/$(TARGET_APISERVER) ./tools/master/master.go

controller:
	$(GO_BUILD) -o ./bin/$(TARGET_CONTROLLER) ./tools/controller/controller.go

scheduler:
	$(GO_BUILD) -o ./bin/$(TARGET_SCHEDULER) ./tools/scheduler/scheduler.go

activer:
	$(GO_BUILD) -o ./bin/$(TARGET_ACTIVER) ./tools/activer/activer.go

kubelet:
	$(GO_BUILD) -o ./bin/$(TARGET_KUBELET) ./tools/kubelet/kubelet.go

proxy:
	$(GO_BUILD) -o ./bin/$(TARGET_PROXY) ./tools/kubeproxy/proxy.go

cadvisor:
	make -C ./third_party/cadvisor
	mv ./third_party/cadvisor/cadvisor ./bin

run_node:
	echo run_node

run_master:
	echo run_master
