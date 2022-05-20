package def

type Function struct {
	Kind         string `yaml:"kind" json:"kind"`
	Name         string
	Function     string `yaml:"function" json:"function"`
	Requirements string `yaml:"requirements" json:"requirements"`
	Version      int    `yaml:"version" json:"version"`
	Image        string `yaml:"image" json:"image"` //yaml文件无此字段
	ServiceName  string //yaml文件无此字段
	PodName      string //yaml文件无此字段
}

type FunctionCache struct {
	Name     string
	Version  int    `yaml:"version" json:"version"`
	Image    string `yaml:"image" json:"image"` //yaml文件无此字段
	Services *ClusterIPSvc
}
