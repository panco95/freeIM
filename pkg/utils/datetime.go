package utils

import "time"

const (
	DefaultTimelayout = "2006-01-02 15:04:05"
)

func StringToTimestamp(datetime string) int64 {
	loc, _ := time.LoadLocation("Local")
	theTime, _ := time.ParseInLocation(DefaultTimelayout, datetime, loc)
	sr := theTime.Unix()
	return sr
}
