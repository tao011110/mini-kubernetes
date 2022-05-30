package def

import (
	"time"
)

type HeartBeat struct {
	NodeID    int       `json:"node_id"`
	TimeStamp time.Time `json:"time_stamp"`
}
