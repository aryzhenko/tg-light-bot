package common

import (
	"fmt"
	"time"
)

func GetNow() time.Time {
	location, _ := time.LoadLocation("Europe/Kiev")
	return time.Now().In(location)
}

func FormatDuration(d time.Duration) string {
	d = d.Round(time.Minute)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	return fmt.Sprintf("%02d год, %02d хв", h, m)
}
