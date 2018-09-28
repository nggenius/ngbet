package toolkit

import (
	"time"
)

const (
	TIME_LAYOUT = "2006-01-02 15:04:05"
)

// UTC时间转本地时间
func UTC2Loc(utc time.Time) time.Time {
	y, mo, d := utc.Date()
	h, mi, s := utc.Clock()
	t := time.Date(y, mo, d, h, mi, s, utc.Nanosecond(), time.Local)
	return t
}

// 两个时间点是否在同一天
func IsSameDay(d1, d2 time.Time) bool {
	return d1.Year() == d2.Year() && d1.Month() == d2.Month() && d1.Day() == d2.Day()
}
