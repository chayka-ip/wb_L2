package server

import (
	"httpserver/calendar"
	"net/url"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParseGetEventQuery(t *testing.T) {
	testCases := []struct {
		name       string
		params     map[string][]string
		wantUserID uint64
		wantDate   time.Time
		wantErr    error
	}{
		{
			name:       "basic",
			params:     map[string][]string{qUserID: {"999"}, qDate: {"2019-12-15"}},
			wantUserID: 999,
			wantDate:   calendar.DateFromYearMonthDay(2019, 12, 15),
		},
		{
			name:    "error parse user id",
			params:  map[string][]string{qUserID: {"-999"}, qDate: {"2019-12-15"}},
			wantErr: errInvalidArgument(qUserID),
		},
		{
			name:    "error parse date",
			params:  map[string][]string{qUserID: {"999"}, qDate: {"2019-13-15"}},
			wantErr: errInvalidArgument(qDate),
		},
		{
			name:    "error no user id",
			params:  map[string][]string{qDate: {"2019-13-15"}},
			wantErr: errNotEnoughArguments(qUserID, qDate),
		},
		{
			name:    "error no date",
			params:  map[string][]string{qUserID: {"999"}},
			wantErr: errNotEnoughArguments(qUserID, qDate),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := parseGetEventQuery(url.Values(tc.params))
			if tc.wantErr != nil {
				assert.Equal(t, tc.wantErr, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tc.wantUserID, res.userID)
			assert.Equal(t, tc.wantDate, res.date)
		})
	}
}

func TestParseEventCRUDStruct(t *testing.T) {
	type ts struct {
		name    string
		data    string
		want    calendar.EventCRUD
		wantErr error
	}

	tc1 := ts{
		name: "basic",
		data: `{"user_id": 11, "date": "2019-05-15",
				"title": "title test", "description":"somedesc"}`,
		want: makeCRUDStruct(11, 0, "2019-05-15", "title test", "somedesc"),
	}
	tc1.want.UID = nil

	tc2 := ts{
		name: "no user id",
		data: `{"date": "2019-05-15",
				"title": "title test", "description":"somedesc"}`,
		want: makeCRUDStruct(0, 0, "2019-05-15", "title test", "somedesc"),
	}
	tc2.want.UID = nil
	tc2.want.UserID = nil

	testCases := []ts{tc1, tc2}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := parseEventCRUDStruct([]byte(tc.data))
			if tc.wantErr != nil {
				assert.Equal(t, tc.wantErr, err)
				return
			}
			assert.NoError(t, err)
			eq := reflect.DeepEqual(tc.want, *got)
			assert.True(t, eq)
		})
	}
}

func TestParseCreateEventRequest(t *testing.T) {
	testCases := []struct {
		name       string
		data       string
		wantUserID uint64
		wantEvent  calendar.Event
		wantErr    error
	}{
		{
			name: "all fields valid",
			data: `{"user_id": 11, "date": "2019-05-15",
				    "title": "title test", "description":"somedesc"}`,
			wantUserID: 11,
			wantEvent: calendar.Event{
				Date:        calendar.DateFromYearMonthDay(2019, 5, 15),
				Title:       "title test",
				Description: "somedesc",
			},
		},
		{
			name: "min setup",
			data: `{"user_id": 11, "date": "2019-05-15",
				    "title": "title test"}`,
			wantUserID: 11,
			wantEvent: calendar.Event{
				Date:  calendar.DateFromYearMonthDay(2019, 5, 15),
				Title: "title test",
			},
		},
		{
			name: "no user_id is err",
			data: `{"date": "2019-05-15",
				    "title": "title test"}`,
			wantErr: errBadRequestData,
		},
		{
			name: "no date is err",
			data: `{"user_id": 11, 
				    "title": "title test"}`,
			wantErr: errBadRequestData,
		},
		{
			name:    "no title is err",
			data:    `{"user_id": 11, "date": "2019-05-15"}`,
			wantErr: errBadRequestData,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			userID, event, err := parseCreateEventRequest([]byte(tc.data))
			if tc.wantErr != nil {
				assert.Equal(t, tc.wantErr, err)
				return
			}
			assert.NoError(t, err)
			eq := reflect.DeepEqual(tc.wantEvent, *event)
			assert.True(t, eq)
			assert.Equal(t, tc.wantUserID, userID)
		})
	}
}

func TestParseUpateEventRequest(t *testing.T) {
	type ts struct {
		name    string
		data    string
		want    calendar.EventCRUD
		wantErr error
	}

	tc1 := ts{
		name: "minimal setup",
		data: `{"user_id": 11, "uid": 55}`,
		want: makeCRUDStruct(11, 55, "", "", ""),
	}
	tc1.want.Date = nil
	tc1.want.Title = nil
	tc1.want.Description = nil

	tc2 := ts{
		name: "update title",
		data: `{"user_id": 11, "uid": 55, "title":"new title"}`,
		want: makeCRUDStruct(11, 55, "", "new title", ""),
	}
	tc2.want.Date = nil
	tc2.want.Description = nil

	testCases := []ts{tc1, tc2}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := parseUpdateEventRequest([]byte(tc.data))
			if tc.wantErr != nil {
				assert.Equal(t, tc.wantErr, err)
				return
			}
			assert.NoError(t, err)
			eq := reflect.DeepEqual(tc.want, *got)
			assert.True(t, eq)
		})
	}
}

func TestParseDeleteEventRequest(t *testing.T) {
	type ts struct {
		name    string
		data    string
		want    deleteEventData
		wantErr error
	}

	tc1 := ts{
		name: "minimal setup",
		data: `{"user_id": 11, "uid": 55}`,
		want: deleteEventData{userID: 11, eventUID: 55},
	}

	tc2 := ts{
		name:    "no user_id is error",
		data:    `{"uid": 55}`,
		wantErr: errBadRequestData,
	}
	tc3 := ts{
		name:    "no uid is error",
		data:    `{"user_id": 11`,
		wantErr: errBadRequestData,
	}

	testCases := []ts{tc1, tc2, tc3}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			data, err := parseDeleteEventRequest([]byte(tc.data))
			if tc.wantErr != nil {
				assert.Equal(t, tc.wantErr, err)
				return
			}
			assert.NoError(t, err)
			eq := reflect.DeepEqual(tc.want, *data)
			assert.True(t, eq)
		})
	}
}

func makeCRUDStruct(userID, UID uint64, date, title, desc string) calendar.EventCRUD {
	return calendar.EventCRUD{
		UserID:      &userID,
		UID:         &UID,
		Date:        &date,
		Title:       &title,
		Description: &desc,
	}
}
