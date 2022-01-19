package models

import "time"

type DateTime struct {
	Year  int
	Month int
	Day   int
}

func (e *DateTime) ToTime() time.Time {
	return time.Date(e.Year, time.Month(e.Month), e.Day, 0, 0, 0, 0, time.UTC)
}
