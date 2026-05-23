package service

import "time"

func combineDateClock(date time.Time, clock string, dayOffset int) time.Time {
	if clock == "" {
		clock = "00:00:00"
	}
	parsed, err := time.ParseInLocation("2006-01-02 15:04:05", date.Format("2006-01-02")+" "+clock, time.Local)
	if err != nil {
		return date
	}
	return parsed.AddDate(0, 0, dayOffset)
}
