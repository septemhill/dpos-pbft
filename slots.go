package main

import (
	"math"
	"time"
)

func beginEpochTime() int64 {
	return time.Date(1987, time.November, 17, 0, 0, 0, 0, time.UTC).Unix()
}

//GetEpochTime get epoch time
func GetEpochTime(times int64) int64 {
	if times == 0 {
		times = time.Now().Unix()
	}

	beginTime := beginEpochTime()
	return times - beginTime
}

//GetTime get time
func GetTime(times int64) int64 {
	return GetEpochTime(times)
}

//GetSlotNumber get slot number
func GetSlotNumber(epochTime int64) int64 {
	if epochTime == 0 {
		epochTime = GetTime(0)
	}
	return int64(math.Floor(float64(epochTime / int64(slotTimeInterval))))
}
