package util

import "time"

const (
	plusIST = 5*time.Hour + 30*time.Minute
)

func GetFormattedTime(t time.Time) string {
	return t.Add(plusIST).Format("Mon Jan 2 15:04:05 IST 2006")
}
