package def

type GPUSlurmConfig struct {
	JobName                  string `yaml:"jobName"`
	Partition                string `yaml:"partition"`
	CpusPerTask              int    `yaml:"cpusPerTask"`
	NtasksPerNode            int    `yaml:"ntasksPerNode"`
	Node                     int    `yaml:"Node"`
	GPU                      int    `yaml:"GPU"`
	Time                     string `yaml:"time"`
	TargetExecutableFileName string `yaml:"targetExecutableFileName"`
}

type GPUJob struct {
	Kind           string         `yaml:"kind"`
	Name           string         `yaml:"name"`
	SourceCodePath string         `yaml:"sourceCodePath"`
	MakefilePath   string         `yaml:"MakefilePath"`
	Slurm          GPUSlurmConfig `yaml:"slurm"`
	ResultPath     string         `yaml:"resultPath"`
	ImageName      string
	PodName        string
}
