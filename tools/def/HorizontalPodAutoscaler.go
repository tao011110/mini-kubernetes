package def

import "time"

type ParsedHorizontalPodAutoscaler struct {
	Name           string    `json:"name"`
	CPUMinValue    float64   `json:"CPUMinValue"`
	CPUMaxValue    float64   `json:"CPUMaxValue"`
	MemoryMinValue string    `json:"memoryMinValue"`
	MemoryMaxValue string    `json:"memoryMaxValue"`
	MinReplicas    int       `json:"minReplicas"`
	MaxReplicas    int       `json:"maxReplicas"`
	PodName        string    `json:"podName"`
	StartTime      time.Time `json:"startTime"`
}
