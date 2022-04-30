package def

import (
	"time"
)

type ResourceUsage struct {
	Time        time.Time `json:"time"`
	CPULoad     int32     `json:"CPULoad"`
	MemoryUsage uint64    `json:"memoryUsage"`
	MemoryTotal uint64    `json:"memoryTotal"`
}

type ResourceUsageSequence struct {
	Sequence []ResourceUsage `json:"sequence"`
}
