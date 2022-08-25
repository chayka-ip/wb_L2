package server

import (
	"encoding/json"
	"httpserver/calendar"
	"net/url"
	"strconv"
	"time"
)

type deleteEventData struct {
	userID   uint64
	eventUID uint64
}

// expected query params:
var (
	qUserID = "user_id"
	qDate   = "date"
)

type getEventQuery struct {
	userID uint64
	date   time.Time
}

func parseGetEventQuery(r url.Values) (*getEventQuery, error) {
	sUserID := r.Get(qUserID)
	sDate := r.Get(qDate)

	if sUserID == "" || sDate == "" {
		return nil, errNotEnoughArguments(qUserID, qDate)
	}

	userID, err := strconv.ParseUint(sUserID, 10, 64)
	if err != nil {
		return nil, errInvalidArgument(qUserID)
	}

	date, err := calendar.DateFromString(sDate)
	if err != nil {
		return nil, errInvalidArgument(qDate)
	}

	q := &getEventQuery{userID: userID, date: date}

	return q, nil
}

func parseEventCRUDStruct(data []byte) (*calendar.EventCRUD, error) {
	e := &calendar.EventCRUD{}
	if err := json.Unmarshal(data, e); err != nil {
		return nil, errBadRequestData
	}
	return e, nil
}

func parseCreateEventRequest(data []byte) (uint64, *calendar.Event, error) {
	badReturn := func() (uint64, *calendar.Event, error) {
		return 0, nil, errBadRequestData
	}
	crud, err := parseEventCRUDStruct(data)
	if err != nil {
		return badReturn()
	}

	ok := crud.HasUserID() && crud.HasTitle() && crud.HasDate()
	if !ok {
		return badReturn()
	}

	e, err := crud.ToCalendarEvent()
	if err != nil {
		return badReturn()
	}

	return *crud.UserID, e, nil
}

func parseUpdateEventRequest(data []byte) (*calendar.EventCRUD, error) {
	badReturn := func() (*calendar.EventCRUD, error) {
		return nil, errBadRequestData
	}

	crud, err := parseEventCRUDStruct(data)
	if err != nil {
		return badReturn()
	}

	ok := crud.HasUserID() && crud.HasUID()
	if !ok {
		return badReturn()
	}

	return crud, nil
}

func parseDeleteEventRequest(data []byte) (*deleteEventData, error) {
	badReturn := func() (*deleteEventData, error) { return nil, errBadRequestData }

	crud, err := parseEventCRUDStruct(data)
	if err != nil {
		return badReturn()
	}

	ok := crud.HasUserID() && crud.HasUID()
	if !ok {
		return badReturn()
	}

	d := &deleteEventData{
		userID:   *crud.UserID,
		eventUID: *crud.UID,
	}

	return d, nil
}
