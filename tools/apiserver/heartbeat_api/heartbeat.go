package heartbeat_api

import (
	"mini-kubernetes/tools/def"
	"time"
)

func ReceiveHeartBeat(beat def.HeartBeat, heartBeatMap map[int]time.Time) {
	heartBeatMap[beat.NodeID] = beat.TimeStamp
}
