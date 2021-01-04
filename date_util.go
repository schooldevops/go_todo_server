package main

import "time"

const TimeFormat = "2006-01-02T15:04:05.999"

func timeToString(time time.Time) string {
	return time.Format(TimeFormat)
}

func stringToTime(timeStr string) time.Time {
	time, _ := time.Parse(TimeFormat, timeStr)
	return time
}
