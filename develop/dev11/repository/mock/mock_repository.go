package mockrepository

import (
	"httpserver/calendar"
	"httpserver/repository"
	"math/rand"
	"time"
)

type eventStorage map[uint64][]calendar.Event

var (
	emptyEvent = calendar.Event{}
)

//MockEventRepository ...
type MockEventRepository struct {
	storage eventStorage
}

//New ...
func New() *MockEventRepository {
	return &MockEventRepository{
		storage: make(eventStorage),
	}
}

//CreateEvent ...
func (r *MockEventRepository) CreateEvent(userID uint64, event calendar.Event) (calendar.Event, error) {
	if event.Title == "" {
		return emptyEvent, repository.ErrRequiredFieldsNotProvided
	}
	event.UID = getRanomUID()
	r.storage[userID] = append(r.storage[userID], event)
	return event, nil
}

//UpdateEvent ...
func (r *MockEventRepository) UpdateEvent(event calendar.EventCRUD) (calendar.Event, error) {
	upd, err := event.ToCalendarEvent()
	if err != nil {
		return emptyEvent, err
	}

	userID := *event.UserID
	ev, has := r.storage[userID]
	if has {
		for i, v := range ev {
			if v.UID == *event.UID {
				if event.HasTitle() {
					if upd.Title == "" {
						return emptyEvent, repository.ErrIncorrectValue
					}
					v.Title = upd.Title
				}
				if event.HasDescription() {
					v.Description = upd.Description
				}
				if event.HasDate() {
					v.Date = upd.Date
				}

				ev[i] = v
				return v, nil
			}
		}
	}
	return emptyEvent, repository.ErrEventDoesNotExists
}

//DeleteEvent ...
func (r *MockEventRepository) DeleteEvent(userID, eventUID uint64) error {
	ev, has := r.storage[userID]
	if has {
		for i, v := range ev {
			if v.UID == eventUID {
				r.storage[userID] = append(ev[:i], ev[i+1:]...)
				return nil
			}
		}
	}
	return repository.ErrEventDoesNotExists
}

//GetEvents ...
func (r *MockEventRepository) GetEvents(userID uint64, date time.Time, dateRange calendar.EventDateRange) ([]calendar.Event, error) {
	ev, has := r.storage[userID]
	if !has {
		return nil, nil
	}

	out := make([]calendar.Event, 0, len(ev))
	tr := calendar.GetTimeRangeFromDateRange(date, dateRange)

	for _, v := range ev {
		if tr.IsTimeInInterval(v.Date) {
			out = append(out, v)
		}
	}

	return out, nil
}

func getRanomUID() uint64 {
	rand.Seed(time.Now().UnixNano())
	return rand.Uint64()
}
