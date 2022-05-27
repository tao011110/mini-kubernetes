package def

import "time"

type ParsedHorizontalPodAutoscaler struct {
	Name           string    `json:"name"`
	CPUMinValue    string    `json:"CPUMinValue"`
	CPUMaxValue    string    `json:"CPUMaxValue"`
	MemoryMinValue string    `json:"memoryMinValue"`
	MemoryMaxValue string    `json:"memoryMaxValue"`
	MinReplicas    int       `json:"minReplicas"`
	MaxReplicas    int       `json:"maxReplicas"`
	PodName        string    `json:"podName"`
	StartTime      time.Time `json:"startTime"`
}
