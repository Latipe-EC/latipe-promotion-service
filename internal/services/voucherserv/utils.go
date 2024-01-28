package voucherserv

import "time"

func ParseStringToTime(dateStr string) time.Time {
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return time.Now()
	}

	return date
}

func ParseDateToString(date time.Time) string {
	layout := "2006-01-02"
	formattedTime := date.Format(layout)
	return formattedTime
}
