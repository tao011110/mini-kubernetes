package def

type PathPair struct {
	Path    string `yaml:"path" json:"path"`
	Service string `yaml:"service" json:"service"`
	Port    uint16 `yaml:"port" json:"port"`
}

type DNS struct {
	Kind  string     `yaml:"kind" json:"kind"`
	Name  string     `yaml:"name" json:"name"`
	Host  string     `yaml:"host" json:"host"`
	Paths []PathPair `yaml:"paths" json:"paths"`
}

type PathPairDetail struct {
	Path    string            `yaml:"path" json:"path"`
	Service ClusterIPInstance `yaml:"service" json:"service"`
	Port    uint16            `yaml:"port" json:"port"`
}

type DNSDetail struct {
	Kind  string           `yaml:"kind" json:"kind"`
	Name  string           `yaml:"name" json:"name"`
	Host  string           `yaml:"host" json:"host"`
	Paths []PathPairDetail `yaml:"paths" json:"paths"`
}

/*TODO: convert DNS to DNSDetail!*/
