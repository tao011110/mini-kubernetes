package yaml

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"mini-kubernetes/tools/def"
	"mini-kubernetes/tools/pod"
	"os"
)

func ReadPodYamlConfig(path string) (*pod.Pod, error) {
	pod_ := &pod.Pod{}
	if f, err := os.Open(path); err != nil {
		return nil, err
	} else {
		err := yaml.NewDecoder(f).Decode(pod_)
		if err != nil {
			return nil, err
		}
		if (*pod_).ApiVersion != "v1" {
			fmt.Println("apiVersion should be v1!")
			return nil, err
		} else if (*pod_).Kind != "Pod" {
			fmt.Println("kind should be Pod!")
			return nil, err
		}
	}
	fmt.Println("pod_: ", pod_)
	return pod_, nil
}

func ReadServiceClusterIPConfig(path string) (*def.ClusterIP, error) {
	service_c := &def.ClusterIP{}
	if f, err := os.Open(path); err != nil {
		return nil, err
	} else {
		err := yaml.NewDecoder(f).Decode(service_c)
		if err != nil {
			return nil, err
		}
		if (*service_c).ApiVersion != "v1" {
			fmt.Println("apiVersion should be v1!")
			return nil, err
		} else if (*service_c).Kind != "Service" {
			fmt.Println("kind should be Pod!")
			return nil, err
		} else if (*service_c).Spec.Type != "ClusterIP" {
			fmt.Println("spec type should be ClusterIP!")
			return nil, err
		}
	}
	fmt.Println("service_c: ", service_c)
	return service_c, nil
}

func ReadServiceNodeportConfig(path string) (*def.Nodeport, error) {
	service_n := &def.Nodeport{}
	if f, err := os.Open(path); err != nil {
		return nil, err
	} else {
		err := yaml.NewDecoder(f).Decode(service_n)
		if err != nil {
			return nil, err
		}
		if (*service_n).ApiVersion != "v1" {
			fmt.Println("apiVersion should be v1!")
			return nil, err
		} else if (*service_n).Kind != "Service" {
			fmt.Println("kind should be Pod!")
			return nil, err
		} else if (*service_n).Spec.Type != "NodePort" {
			fmt.Println("spec type should be NodePort!")
			return nil, err
		}
	}
	fmt.Println("service_n: ", service_n)
	return service_n, nil
}

func ReadDNSConfig(path string) (*def.DNS, error) {
	dns_ := &def.DNS{}
	if f, err := os.Open(path); err != nil {
		return nil, err
	} else {
		err := yaml.NewDecoder(f).Decode(dns_)
		if err != nil {
			return nil, err
		}
		if (*dns_).Kind != "DNS" {
			fmt.Println("kind should be DNS!")
			return nil, err
		}
	}
	fmt.Println("dns_: ", dns_)
	return dns_, nil
}
