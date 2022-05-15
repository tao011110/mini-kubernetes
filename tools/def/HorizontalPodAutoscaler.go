package def

type ParsedHorizontalPodAutoscaler struct {
	Name           string
	CPUMinValue    float64
	CPUMaxValue    float64
	MemoryMinValue int64
	MemoryMaxValue int64
	MinReplicas    int
	MaxReplicas    int
	PodName        string
}
