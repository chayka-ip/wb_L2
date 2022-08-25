package diskvrepository

import (
	"httpserver/calendar"
	"httpserver/repository"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func newTestRepository() *DiskvEventRepository {
	cacheSize := 1024 * 1024
	dir := "db_test"
	r := New(dir, uint64(cacheSize))
	r.Erase()
	return r
}

func TestDiskvRepoCreateEvent(t *testing.T) {
	r := newTestRepository()

	testCases := []struct {
		name              string
		userID            uint64
		event             calendar.Event
		expectedErr       error
		expectedEventList []calendar.Event
	}{
		{
			name:   "basic",
			userID: 1,
			event: calendar.Event{
				Title:       "test",
				Description: "desc",
				Date:        time.Now(),
			},
		},
		{
			name:   "no title is error",
			userID: 1,
			event: calendar.Event{
				Description: "desc",
				Date:        time.Now(),
			},
			expectedErr: repository.ErrRequiredFieldsNotProvided,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			defer r.Erase()
			e, err := r.CreateEvent(tc.userID, tc.event)
			if tc.expectedErr != nil {
				assert.Equal(t, tc.expectedErr, err)
				return
			}
			assert.Equal(t, tc.event.Title, e.Title)
			assert.Equal(t, tc.event.Description, e.Description)
			assert.Equal(t, tc.event.Date, e.Date)
		})
	}

}
func TestDiskvRepoUpdateEvent(t *testing.T) {
	r := newTestRepository()

	var userID uint64 = 1
	event := calendar.Event{
		Title:       "title",
		Description: "description",
		Date:        time.Now(),
	}

	event, err := r.CreateEvent(userID, event)
	assert.NoError(t, err)
	var eventUID uint64 = event.UID

	u := calendar.EventCRUD{
		UserID:      &userID,
		UID:         &eventUID,
		Description: &event.Description,
	}

	{ // test update existing event for wrong user
		wrongUserID := userID + 1
		u.UserID = &wrongUserID

		_, err = r.UpdateEvent(u)
		assert.Equal(t, repository.ErrEventDoesNotExists, err)
	}

	{ // test update not existing event for existing user
		wrongUID := eventUID + 1
		u.UserID = &userID
		u.UID = &wrongUID

		_, err = r.UpdateEvent(u)
		assert.Equal(t, repository.ErrEventDoesNotExists, err)
	}

	testCases := []struct {
		name          string
		event         calendar.Event
		update        calendar.EventCRUD
		expectedEvent calendar.Event
		expectedErr   error
	}{
		{
			name: "update description",
			event: calendar.Event{
				Title:       "Title",
				Description: "Description",
				Date:        calendar.DateNow(),
			},
			update: calendar.EventCRUD{
				Description: repository.StringPtr("New Description"),
			},
			expectedEvent: calendar.Event{
				Title:       "Title",
				Description: "New Description",
				Date:        calendar.DateNow(),
			},
		},
		{
			name: "update date",
			event: calendar.Event{
				Title:       "Title",
				Description: "Description",
				Date:        calendar.DateNow(),
			},
			update: calendar.EventCRUD{
				Date: repository.StringPtr("2000-05-13"),
			},
			expectedEvent: calendar.Event{
				Title:       "Title",
				Description: "Description",
				Date:        calendar.DateFromYearMonthDay(2000, 5, 13),
			},
		},
		{
			name: "update title",
			event: calendar.Event{
				Title:       "Title",
				Description: "Description",
				Date:        calendar.DateNow(),
			},
			update: calendar.EventCRUD{
				Title: repository.StringPtr("New Title"),
			},
			expectedEvent: calendar.Event{
				Title:       "New Title",
				Description: "Description",
				Date:        calendar.DateNow(),
			},
		},
		{
			name: "update title to empty is error",
			event: calendar.Event{
				Title:       "Title",
				Description: "Description",
				Date:        calendar.DateNow(),
			},
			update: calendar.EventCRUD{
				Title: repository.StringPtr(""),
			},
			expectedErr: repository.ErrIncorrectValue,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			defer r.Erase()

			var userID uint64 = 1
			e, _ := r.CreateEvent(userID, tc.event)

			// set mandatory fields
			tc.update.UserID = &userID
			tc.update.UID = &e.UID

			eventUpd, err := r.UpdateEvent(tc.update)
			if tc.expectedErr != nil {
				assert.Equal(t, tc.expectedErr, err)
				return
			}

			tc.expectedEvent.UID = e.UID
			eq := reflect.DeepEqual(tc.expectedEvent, eventUpd)
			assert.True(t, eq)
		})
	}
}
func TestDiskvRepoDeleteEvent(t *testing.T) {
	r := newTestRepository()
	defer r.Erase()

	var userID uint64 = 1
	e := calendar.Event{
		Title:       "t",
		Description: "desc",
		Date:        time.Now(),
	}

	e, err := r.CreateEvent(userID, e)
	assert.NoError(t, err)
	var eventUID uint64 = e.UID

	{ // delete not existing event for existing user
		err = r.DeleteEvent(userID, eventUID+1)
		assert.Equal(t, repository.ErrEventDoesNotExists, err)
	}

	{ // delete existing event but for wrong user
		err = r.DeleteEvent(userID+1, eventUID)
		assert.Equal(t, repository.ErrEventDoesNotExists, err)
	}

	{ // delete existing event for existing user
		err = r.DeleteEvent(userID, eventUID)
		assert.NoError(t, err)
	}

	{ // delete removed event to get error
		err = r.DeleteEvent(userID, eventUID)
		assert.Equal(t, repository.ErrEventDoesNotExists, err)
	}

}
func TestDiskvRepoGetEvents(t *testing.T) {
	me := repository.MakeEventFromYearMonthDay

	testCases := []struct {
		name           string
		date           time.Time
		dateRange      calendar.EventDateRange
		eventsToCreate []calendar.Event
		expectedEvents []calendar.Event
		expectedErr    error
	}{
		{
			name:      "get day events",
			dateRange: calendar.EventDAY,
			date:      calendar.DateFromYearMonthDay(2020, 04, 15),
			eventsToCreate: []calendar.Event{
				me(2020, 04, 15), me(2020, 04, 15),
				me(2020, 04, 20), me(2020, 05, 15)},
			expectedEvents: []calendar.Event{
				me(2020, 04, 15), me(2020, 04, 15),
			},
		},
		{
			name:      "get week events",
			dateRange: calendar.EventWEEK,
			date:      calendar.DateFromYearMonthDay(2020, 04, 15),
			eventsToCreate: []calendar.Event{
				me(2020, 04, 13), me(2020, 04, 17),
				me(2020, 04, 20), me(2020, 05, 15)},
			expectedEvents: []calendar.Event{
				me(2020, 04, 13), me(2020, 04, 17),
			},
		},
		{
			name:      "get month events",
			dateRange: calendar.EventMONTH,
			date:      calendar.DateFromYearMonthDay(2020, 04, 15),
			eventsToCreate: []calendar.Event{
				me(2020, 04, 1), me(2020, 04, 30),
				me(2020, 02, 20), me(2020, 05, 15)},
			expectedEvents: []calendar.Event{
				me(2020, 04, 1), me(2020, 04, 30),
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := newTestRepository()
			defer r.Erase()

			var userID uint64 = 1
			for _, v := range tc.eventsToCreate {
				r.CreateEvent(userID, v)
			}
			result, err := r.GetEvents(userID, tc.date, tc.dateRange)
			if tc.expectedErr != nil {
				assert.Equal(t, tc.expectedErr, err)
				return
			}

			eq := calendar.AreEqualEventSlicesByDateAndCount(tc.expectedEvents, result)
			assert.True(t, eq)
		})
	}
}
