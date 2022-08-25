package calendar

import (
	"time"
)

//EventDateRange determines type of event selection
type EventDateRange uint8

const (
	//EventDAY is an event for specified day
	EventDAY EventDateRange = 0
	//EventWEEK is an event on the week to which specific day belongs
	EventWEEK EventDateRange = 1
	//EventMONTH is an event on the month to which specific day belongs
	EventMONTH EventDateRange = 2

	//TimeLayout represents date format (YYYY-MM-DD)
	TimeLayout = "2006-01-02"
)

//Event represents event
type Event struct {
	UID         uint64    `json:"uid,omitempty"`
	Date        time.Time `json:"date,omitempty"`
	Title       string    `json:"title,omitempty"`
	Description string    `json:"description,omitempty"`
}

//EventRepository is an interface to work with calendar events
type EventRepository interface {
	CreateEvent(userID uint64, event Event) (Event, error)
	UpdateEvent(event EventCRUD) (Event, error)
	DeleteEvent(userID, eventUID uint64) error
	GetEvents(userID uint64, date time.Time, dateRange EventDateRange) ([]Event, error)
}
