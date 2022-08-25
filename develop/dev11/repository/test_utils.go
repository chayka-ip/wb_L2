package repository

import (
	"httpserver/calendar"
	"time"
)

//Uint64Ptr ...
func Uint64Ptr(i uint64) *uint64 {
	return &i
}

//StringPtr ...
func StringPtr(s string) *string {
	return &s
}

//MakeEventFromYearMonthDay ...
func MakeEventFromYearMonthDay(year int, month time.Month, day int) calendar.Event {
	return calendar.Event{
		Title: "T",
		Date:  calendar.DateFromYearMonthDay(year, month, day),
	}
}
