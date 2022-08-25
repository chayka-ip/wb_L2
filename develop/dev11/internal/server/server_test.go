package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"httpserver/calendar"
	mockrepository "httpserver/repository/mock"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

type postTestCase struct {
	name     string
	data     string
	wantCode int
}

func newTestServer() *Server {
	config := NewDefaultConfig()
	repo := mockrepository.New()
	logger := logrus.New()
	return New(config, repo, logger)
}

func TestHandleCreateEvent(t *testing.T) {
	s := newTestServer()

	testCases := []postTestCase{
		{
			name: "all fields valid",
			data: `{"user_id": 11, "date": "2019-05-15",
				    "title": "title test", "description":"somedesc"}`,
			wantCode: http.StatusCreated,
		},
		{
			name: "min setup",
			data: `{"user_id": 11, "date": "2019-05-15",
				    "title": "title test"}`,
			wantCode: http.StatusCreated,
		},
		{
			name: "no user_id is err",
			data: `{"date": "2019-05-15",
				    "title": "title test"}`,
			wantCode: http.StatusBadRequest,
		},
		{
			name: "no date is err",
			data: `{"user_id": 11,
				    "title": "title test"}`,
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "no title is err",
			data:     `{"user_id": 11, "date": "2019-05-15"}`,
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "invalid date format is err",
			data:     `{"user_id": 11, "date": "2019-15-05"}`,
			wantCode: http.StatusBadRequest,
		},
		{
			name: "invalid user id format is err",
			data: `{"user_id": "11", "date": "2019-05-15",
				    "title": "title test"}`,
			wantCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := getResponse(s, http.MethodPost, apiCreateEvent, []byte(tc.data))
			assert.Equal(t, tc.wantCode, rec.Code)
		})
	}
}
func TestHandleUpdateEvent(t *testing.T) {
	s := newTestServer()

	// setup test data
	var userID uint64 = 1
	event := calendar.Event{Title: "title", Description: "descrption"}
	event, err := s.eventRepo.CreateEvent(userID, event)
	eventUID := event.UID
	assert.NoError(t, err)

	testCases := []postTestCase{
		{
			name:     "update not existing event for existing user is err",
			data:     fmt.Sprintf(`{"user_id": %d, "uid": %d, "title":"new title"}`, userID, eventUID+1),
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "update existing event for wrong user is err",
			data:     fmt.Sprintf(`{"user_id": %d, "uid": %d, "title":"new title"}`, userID+1, eventUID),
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "update description",
			data:     fmt.Sprintf(`{"user_id": %d, "uid": %d, "descrption":"new text"}`, userID, eventUID),
			wantCode: http.StatusOK,
		},
		{
			name:     "update date",
			data:     fmt.Sprintf(`{"user_id": %d, "uid": %d, "date":"2000-05-13"}`, userID, eventUID),
			wantCode: http.StatusOK,
		},
		{
			name:     "update title",
			data:     fmt.Sprintf(`{"user_id": %d, "uid": %d, "title":"new title"}`, userID, eventUID),
			wantCode: http.StatusOK,
		},
		{
			name:     "update title to empty is error",
			data:     fmt.Sprintf(`{"user_id": %d, "uid": %d, "title":""}`, userID, eventUID),
			wantCode: http.StatusBadRequest,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := getResponse(s, http.MethodPost, apiUpdateEvent, []byte(tc.data))
			assert.Equal(t, tc.wantCode, rec.Code)
		})
	}
}

func TestHandleDeleteEvent(t *testing.T) {
	s := newTestServer()

	// setup test data
	var userID uint64 = 1
	event := calendar.Event{Title: "title", Description: "descrption"}
	event, err := s.eventRepo.CreateEvent(userID, event)
	eventUID := event.UID
	assert.NoError(t, err)

	testCases := []postTestCase{
		{
			name:     "delete not existing event for existing user is err",
			data:     fmt.Sprintf(`{"user_id": %d, "uid": %d}`, userID, eventUID+1),
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "delete existing event but for wrong user is err",
			data:     fmt.Sprintf(`{"user_id": %d, "uid": %d}`, userID+1, eventUID),
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "delete existing event for existing user",
			data:     fmt.Sprintf(`{"user_id": %d, "uid": %d}`, userID, eventUID),
			wantCode: http.StatusNoContent,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := getResponse(s, http.MethodPost, apiDeleteEvent, []byte(tc.data))
			assert.Equal(t, tc.wantCode, rec.Code)
		})
	}
}

func TestHandleGetEvents(t *testing.T) {
	me := makeEventFromYearMonthDay
	type res struct {
		Result []calendar.Event `json:"result"`
	}

	testCases := []struct {
		name              string
		endpoint          string
		createEventUserID uint64
		requestUserID     uint64
		date              string
		eventsToCreate    []calendar.Event
		expectedEvents    []calendar.Event
		wantCode          int
		shouldFail        bool
	}{
		{
			name:              "get day events",
			endpoint:          apiEventsForDay,
			createEventUserID: 1,
			requestUserID:     1,
			date:              "2020-04-15",
			eventsToCreate: []calendar.Event{
				me(2020, 04, 15), me(2020, 04, 15),
				me(2020, 04, 20), me(2020, 05, 15)},
			expectedEvents: []calendar.Event{
				me(2020, 04, 15), me(2020, 04, 15),
			},
			wantCode: http.StatusOK,
		},
		{
			name:              "get week events",
			endpoint:          apiEventsForWeek,
			createEventUserID: 1,
			requestUserID:     1,
			date:              "2020-04-15",
			eventsToCreate: []calendar.Event{
				me(2020, 04, 13), me(2020, 04, 17),
				me(2020, 04, 20), me(2020, 05, 15)},
			expectedEvents: []calendar.Event{
				me(2020, 04, 13), me(2020, 04, 17),
			},
			wantCode: http.StatusOK,
		},
		{
			name:              "get month events",
			endpoint:          apiEventsForMonth,
			createEventUserID: 1,
			requestUserID:     1,
			date:              "2020-04-15",
			eventsToCreate: []calendar.Event{
				me(2020, 04, 1), me(2020, 04, 30),
				me(2020, 02, 20), me(2020, 05, 15)},
			expectedEvents: []calendar.Event{
				me(2020, 04, 1), me(2020, 04, 30),
			},
			wantCode: http.StatusOK,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			s := newTestServer()
			for _, v := range tc.eventsToCreate {
				_, err := s.eventRepo.CreateEvent(tc.createEventUserID, v)
				assert.NoError(t, err)
			}
			url := fmt.Sprintf("%s?%s=%d&%s=%s", tc.endpoint, qUserID, tc.requestUserID, qDate, tc.date)

			rec := getResponse(s, http.MethodGet, url, nil)

			assert.Equal(t, tc.wantCode, rec.Code)
			if tc.shouldFail {
				return
			}

			got := res{}
			err := json.Unmarshal(rec.Body.Bytes(), &got)
			assert.NoError(t, err)

			eq := calendar.AreEqualEventSlicesByDateAndCount(tc.expectedEvents, got.Result)
			assert.True(t, eq)
		})
	}
}

func getResponse(s *Server, method, url string, data []byte) *httptest.ResponseRecorder {
	b := &bytes.Buffer{}
	if data != nil {
		b.Write(data)
	}
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(method, url, b)
	s.ServeHTTP(rec, req)
	return rec
}

func makeEventFromYearMonthDay(year int, month time.Month, day int) calendar.Event {
	return calendar.Event{
		Title: "T",
		Date:  calendar.DateFromYearMonthDay(year, month, day),
	}
}
