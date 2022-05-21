package gateway

import (
	"fmt"
	"github.com/ghodss/yaml"
	"mini-kubernetes/tools/def"
)

type Server struct {
	Port int `json:"port"`
}

type Application struct {
	Name string `json:"name"`
}
type Spring struct {
	Application Application `json:"application"`
}

type PathAndUrl struct {
	Path string `json:"path"`
	Url  string `json:"url"`
}

type Zuul struct {
	Routes map[string]PathAndUrl `json:"routes"`
}

type ApplicationYaml struct {
	Server Server `json:"server"`
	Spring Spring `json:"spring"`
	Zuul   Zuul   `json:"zuul"`
}

func GenerateApplicationYaml(dns def.DNSDetail) string {
	application := ApplicationYaml{
		Zuul: Zuul{
			Routes: map[string]PathAndUrl{},
		},
		Server: Server{
			Port: 80,
		},
		Spring: Spring{
			Application: Application{
				Name: "zuul",
			},
		},
	}
	for index, path := range dns.Paths {
		application.Zuul.Routes[fmt.Sprintf("route%d", index)] = PathAndUrl{
			Path: fmt.Sprintf("%s/**", path.Path),
			Url:  fmt.Sprintf("http://%s:%d", path.Service.Spec.ClusterIP, path.Port),
		}
	}
	bytes, _ := yaml.Marshal(application)
	str := string(bytes)
	return str
}
