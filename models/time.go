package models

import "time"

type DateTime struct {
	Year  int32
	Month int32
	Day   int32
}

func (e *DateTime) ToTime() time.Time {
	return time.Date(int(e.Year), time.Month(e.Month), int(e.Day), 0, 0, 0, 0, time.UTC)
}
