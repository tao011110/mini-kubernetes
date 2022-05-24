package def

type Activity struct {
	Kind         string `yaml:"kind" json:"kind"`
	Name         string `yaml:"name" json:"name"`
	Function     string `yaml:"function" json:"function"`
	Requirements string `yaml:"requirements" json:"requirements"`
	Version      int    `yaml:"version" json:"version"`
}
