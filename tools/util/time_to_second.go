package util

import "time"

func TimeToSecond(time time.Time) int {
	return ((time.Day()*24+time.Hour())*60+time.Minute())*60 + time.Second()
}
