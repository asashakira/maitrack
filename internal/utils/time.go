package utils

import (
	"time"
)

// layout is "2006-01-02 15:04"
func StringToUTCTime(value string) (time.Time, error) {
	// parse time as JST time
	jst, _ := time.LoadLocation("Asia/Tokyo")
	jstTime, err := time.ParseInLocation("2006-01-02 15:04", value, jst)
	if err != nil {
		return time.Time{}, err
	}

	// convert to UTC time
	utcTime := jstTime.UTC()

	return utcTime, nil
}
