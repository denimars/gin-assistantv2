package code

func TimeHelper() string {
	return `
package helper

import (
	"fmt"
	"time"
)

func UtcTime() time.Time {
	return time.Now().UTC()
}

func StringToDate(dateStr string) (time.Time, error) {
	layout := "2006-01-02"
	return time.Parse(layout, dateStr)
}

func FormatDateToString(date time.Time) string {
	return date.Format("2006-01-02")
}

func TodayDateString() string {
	return time.Now().Format("2006-01-02")
}

func ParseAndFormatDate(dateStr string) (string, error) {
	date, err := StringToDate(dateStr)
	if err != nil {
		return "", fmt.Errorf("invalid date format: %v", err)
	}
	return FormatDateToString(date), nil
}
	`
}
