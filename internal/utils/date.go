package utils

import (
	"time"
)

// Common date formats
const (
	DateFormatYYYYMMDD = "2006-01-02"
	DateFormatDDMMYYYY = "02-01-2006"
	DateTimeFormat     = "2006-01-02 15:04:05"
	TimeFormat         = "15:04:05"
)

// FormatDate formats a time.Time to a string using the specified layout
func FormatDate(t time.Time, layout string) string {
	return t.Format(layout)
}

// ParseDate parses a date string using the specified layout
func ParseDate(dateStr, layout string) (time.Time, error) {
	return time.Parse(layout, dateStr)
}

// GetStartOfDay returns the start time of the given date (00:00:00)
func GetStartOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

// GetEndOfDay returns the end time of the given date (23:59:59)
func GetEndOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 999999999, t.Location())
}

// IsWeekend checks if the given date falls on a weekend
func IsWeekend(t time.Time) bool {
	weekday := t.Weekday()
	return weekday == time.Saturday || weekday == time.Sunday
}

// AddWorkdays adds the specified number of workdays (excluding weekends) to the given date
func AddWorkdays(t time.Time, days int) time.Time {
	current := t
	if days > 0 {
		for days > 0 {
			current = current.AddDate(0, 0, 1)
			if !IsWeekend(current) {
				days--
			}
		}
	} else {
		for days < 0 {
			current = current.AddDate(0, 0, -1)
			if !IsWeekend(current) {
				days++
			}
		}
	}
	return current
}
