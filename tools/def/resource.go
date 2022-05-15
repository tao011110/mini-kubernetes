package def

import (
	"time"
)

type ResourceUsage struct {
	Time        time.Time `json:"time"`
	CPULoad     int32     `json:"CPULoad"` //*1000
	CPUNum      int       `json:"cpu-num"`
	MemoryUsage uint64    `json:"memoryUsage"`
	MemoryTotal uint64    `json:"memoryTotal"`
	Valid       bool      `json:"vailid"`
}

//type ResourceUsageSequence struct {
//	Sequence []ResourceUsage `json:"sequence"`
//}
