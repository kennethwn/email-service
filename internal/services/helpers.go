package services

import (
	"errors"
	"time"
)

func GetDifferenceTime(now string, last string) (time.Duration, error) {
	parsedNow, err := ConvertStringToDateISO8601(now)
	if err != nil {
		return 0, err
	}

	parsedLastTime, err := ConvertStringToDateISO8601(last)
	if err != nil {
		return 0, err
	}

	diff := parsedNow.Sub(parsedLastTime)
	if diff < 0 {
		diff = -diff
	}

	return diff, nil
}

func ConvertStringToDateISO8601(dateStr string) (time.Time, error) {
	date, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		return time.Now(), errors.New("Can not convert String to DATE with Format ISO8601. Your Input is " + dateStr)
	}
	return date, err
}
