package util

import (
	"os"
	"strings"
	"time"
)

var Tz, _ = time.LoadLocation("Europe/Bratislava")

//Returns true if the file on path a was modified later than the file on path b.
//If an error is encountered, returns false and the error.
func IsNewer(a, b string) (bool, error) {
	fia, err := os.Stat(a)
	if err != nil {
		return false, err
	}

	fib, err := os.Stat(b)
	if err != nil {
		return false, err
	}

	return fia.ModTime().After(fib.ModTime()), nil
}

//NormalizeDate strips the time part from time.Date leaving only
//year, month and day.
func NormalizeDate(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, Tz)
}

//NormalizeTime strips the date part from time.Date leaving only
//hours, minutes, seconds and nanoseconds.
//If forceTZ is true, its location will be set to util.Tz,
//if false, it will be left as is.
func NormalizeTime(t time.Time) time.Time {
	return time.Date(0, 1, 1, t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), Tz)
}

func GetFileType(filename string) string {
	parts := strings.Split(filename, ".")
	if len(parts) < 2 {
		return ""
	}

	return parts[len(parts)-1]
}
