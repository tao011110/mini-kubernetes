package def

type AutoMeta struct {
	Name string `yaml:"name" json:"name"`
}

type CPUDetail struct {
	TargetMinValue string `yaml:"targetMinValue" json:"target_min_value"`
	TargetMaxValue string `yaml:"targetMaxValue" json:"target_max_value"`
}

type MemoryDetail struct {
	TargetMinValue string `yaml:"targetMinValue" json:"target_min_value"`
	TargetMaxValue string `yaml:"targetMaxValue" json:"target_max_value"`
}

type Metrics struct {
	CPU    CPUDetail    `yaml:"CPU" json:"CPU"`
	Memory MemoryDetail `yaml:"memory" json:"memory"`
}

type AutoTemp struct {
	Metadata TempMeta `yaml:"metadata" json:"metadata"`
	Spec     TempSpec `yaml:"spec" json:"spec"`
}

type AutoSpec struct {
	MinReplicas uint64   `yaml:"minReplicas" json:"min_replicas"`
	MaxReplicas uint64   `yaml:"maxReplicas" json:"max_replicas"`
	Metrics     Metrics  `yaml:"metrics" json:"metrics"`
	Template    AutoTemp `yaml:"template" json:"template"`
}

type Autoscaler struct {
	ApiVersion string   `yaml:"apiVersion" json:"api_version"`
	Kind       string   `yaml:"kind" json:"kind"`
	Metadata   AutoMeta `yaml:"metadata" json:"metadata"`
	Spec       AutoSpec `yaml:"spec" json:"spec"`
}
