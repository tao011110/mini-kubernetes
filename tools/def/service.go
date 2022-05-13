package def

type Selector struct {
	Name string `yaml:"name" json:"name"`
}

type Spec struct {
	Type      string     `yaml:"type" json:"type"`
	ClusterIP string     `yaml:"clusterIP" json:"clusterIP"`
	Ports     []PortPair `yaml:"ports" json:"ports"`
	Selector  Selector   `yaml:"selector" json:"selector"`
}

type Labels struct {
	Name string `yaml:"name" json:"name"`
}

type Meta struct {
	Name string `yaml:"name" json:"name"`
	//Labels Labels `yaml:"labels" json:"labels"`
}

type PortPair struct {
	Port       uint16 `yaml:"port" json:"port"`
	TargetPort string `yaml:"targetPort" json:"targetPort"`
	Protocol   string `yaml:"protocol" json:"protocol"`
	NodePort   uint16 `yaml:"nodePort" json:"nodePort"`
}

type NodePortSvc struct {
	ApiVersion string `yaml:"apiVersion" json:"api_version"`
	Kind       string `yaml:"kind" json:"kind"`
	Metadata   Meta   `yaml:"metadata" json:"metadata"`
	Spec       Spec   `yaml:"spec" json:"spec"`
}

type ClusterIPSvc struct {
	ApiVersion string `yaml:"apiVersion" json:"api_version"`
	Kind       string `yaml:"kind" json:"kind"`
	Metadata   Meta   `yaml:"metadata" json:"metadata"`
	Spec       Spec   `yaml:"spec" json:"spec"`
}

type PortsBindings struct {
	Ports     PortPair `yaml:"ports" json:"ports"`
	Endpoints []string `yaml:"endpoints" json:"endpoints"`
}

type Service struct {
	Name          string          `yaml:"name" json:"name"`
	Selector      Selector        `yaml:"selector" json:"selector"`
	Type          string          `yaml:"type" json:"type"`
	IP            string          `yaml:"IP" json:"IP"`
	PortsBindings []PortsBindings `yaml:"portsBindings" json:"portsBindings"`
}
