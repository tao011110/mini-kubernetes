package def

type Selector struct {
	Name string `yaml:"name" json:"name"`
}

type Spec struct {
	Type     string     `yaml:"type" json:"type"`
	Ports    []PortPair `yaml:"ports" json:"ports"`
	Selector Selector   `yaml:"selector" json:"selector"`
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
	NodePort   uint16 `yaml:"nodePort" json:"nodePort"`
	Protocol   string `yaml:"protocol" json:"protocol"`
}

type Nodeport struct {
	ApiVersion string `yaml:"apiVersion" json:"api_version"`
	Kind       string `yaml:"kind" json:"kind"`
	Metadata   Meta   `yaml:"metadata" json:"metadata"`
	Spec       Spec   `yaml:"spec" json:"spec"`
}

type ClusterIP struct {
	ApiVersion string `yaml:"apiVersion" json:"api_version"`
	Kind       string `yaml:"kind" json:"kind"`
	Metadata   Meta   `yaml:"metadata" json:"metadata"`
	Spec       Spec   `yaml:"spec" json:"spec"`
}
