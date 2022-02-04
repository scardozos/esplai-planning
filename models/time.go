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

// Calculates number of weeks for which to know their respective state in the future
// Takes into account startDate, the requested date and the list of days in which state won't change
func CalcWeekNumNoWeeks(startDate time.Time, requestedDate time.Time, nonWeeks []time.Time) int {
	// Change requestedDate to that week's Monday in order to preserve state
	requestedDate = ChangeWeekDay(requestedDate, time.Monday)

	// Compute the amount of days from startDate
	days := requestedDate.Sub(startDate).Hours() / 24
	// Compute the amount of weeks from startDate
	// taking into account the amount of days
	weeks := int(days / 7)

	// initialize sub var, which accounts for the total number
	// of static weeks that have occurred after startDate
	// and before the requestedDate
	var sub int
	for _, time := range nonWeeks {
		if time.After(startDate) && time.Before(requestedDate) {
			sub += 1
		}
	}

	// returns the number of weeks that passed from startDate
	// minus the amount of static weeks that passed in between
	return weeks - sub
}

// Change week day "from" `Time.Time` to a given weekday "to" `time.Weekday`
// Returns time.Time
func ChangeWeekDay(from time.Time, to time.Weekday) time.Time {
	if currentWeekDay := int(from.Weekday()); currentWeekDay != int(to) {
		sub := currentWeekDay
		if currentWeekDay == 0 {
			sub = 7
		}
		return from.AddDate(0, 0, int(to)-sub)
	}
	return from
}
