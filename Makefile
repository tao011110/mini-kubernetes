GO_CMD=go
GO_BUILD=$(GO_CMD) build
GO_CLEAN=$(GO_CMD) clean
GO_TEST=$(GO_CMD) test
GO_GET=$(GO_CMD) get


GO_TEST_PATH= './test/apiserver_test' \
			'./test/application_yaml_config_test' \
			'./test/django_test' \
			'./test/dns_test' \
			'./test/docker_test' \
			'./test/echo_test' \
			'./test/etcd_test' \
			'./test/image_factory_test' \
			'./test/kubectl_test' \
			'./test/kubeproxy_test' \
			'./test/master_test' \
			'./test/net_utils_test' \
			'./test/resource_test' \
			'./test/slurmGenrator_test' \
			'./test/yaml_test'
.PHONY:test

all: test kubectl apiServer controller scheduler activer kubelet

test:
	for path in ${GO_TEST_PATH}; do \
	$(GO_TEST) $$path -v ; \
	done

kubectl:

apiServer:

controller:

scheduler:
activer:
kubelet:

