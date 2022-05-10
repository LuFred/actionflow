package timeutil

import "time"

func GetTimeMillisecond() int64 {
	return time.Now().UnixNano() / 1e6
}

func GetAddTimeMillisecond(years, month, day int) int64 {
	return time.Now().AddDate(years, month, day).UnixNano() / 1e6
}
