package diskvrepository

import (
	"encoding/json"
	"fmt"
	"httpserver/calendar"
	"httpserver/repository"
	"math/rand"
	"time"

	"github.com/peterbourgon/diskv/v3"
)

var (
	emptyEvent = calendar.Event{}
)

type storageValue map[uint64]calendar.Event

//DiskvEventRepository ...
type DiskvEventRepository struct {
	db *diskv.Diskv
}

//New ...
func New(storageDir string, cacheSize uint64) *DiskvEventRepository {
	// transform function: put all the data files into the base dir.
	flatTransform := func(s string) []string { return []string{} }

	db := diskv.New(diskv.Options{
		BasePath:     storageDir,
		Transform:    flatTransform,
		CacheSizeMax: cacheSize,
	})

	return &DiskvEventRepository{
		db: db,
	}
}

func (r *DiskvEventRepository) read(userID uint64) (storageValue, error) {
	k := uitoa(userID)
	value := make(storageValue)

	if !r.db.Has(k) {
		return value, nil
	}

	b, err := r.db.Read(k)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(b, &value)
	if err != nil {
		return nil, err
	}
	return value, nil
}

func (r *DiskvEventRepository) write(userID uint64, value storageValue) error {
	k := uitoa(userID)
	val, err := json.Marshal(value)
	if err != nil {
		return nil
	}
	return r.db.Write(k, val)
}

//Erase erases all data from the storage
func (r *DiskvEventRepository) Erase() error {
	return r.db.EraseAll()
}

//CreateEvent ...
func (r *DiskvEventRepository) CreateEvent(userID uint64, event calendar.Event) (calendar.Event, error) {
	if event.Title == "" {
		return emptyEvent, repository.ErrRequiredFieldsNotProvided
	}
	uid := getRanomUID()
	event.UID = uid

	value, err := r.read(userID)
	if err != nil {
		return emptyEvent, err
	}

	value[uid] = event
	r.write(userID, value)

	return event, nil
}

//UpdateEvent ...
func (r *DiskvEventRepository) UpdateEvent(event calendar.EventCRUD) (calendar.Event, error) {
	upd, err := event.ToCalendarEvent()
	if err != nil {
		return emptyEvent, err
	}

	userID := *event.UserID
	eventUID := upd.UID
	value, err := r.read(userID)
	if err != nil {
		return emptyEvent, err
	}

	ev, has := value[eventUID]
	if has {
		if event.HasTitle() {
			if upd.Title == "" {
				return emptyEvent, repository.ErrIncorrectValue
			}
			ev.Title = upd.Title
		}
		if event.HasDescription() {
			ev.Description = upd.Description
		}
		if event.HasDate() {
			ev.Date = upd.Date
		}
		value[eventUID] = ev
		r.write(userID, value)
		return ev, nil
	}
	return emptyEvent, repository.ErrEventDoesNotExists
}

//DeleteEvent ...
func (r *DiskvEventRepository) DeleteEvent(userID, eventUID uint64) error {
	value, err := r.read(userID)
	if err != nil {
		return err
	}

	_, has := value[eventUID]
	if has {
		delete(value, eventUID)
		r.write(userID, value)
		return nil
	}
	return repository.ErrEventDoesNotExists
}

//GetEvents ...
func (r *DiskvEventRepository) GetEvents(userID uint64, date time.Time, dateRange calendar.EventDateRange) ([]calendar.Event, error) {
	value, err := r.read(userID)
	if err != nil {
		return nil, err
	}
	if len(value) == 0 {
		return nil, nil
	}

	out := make([]calendar.Event, 0, len(value))
	tr := calendar.GetTimeRangeFromDateRange(date, dateRange)

	for _, v := range value {
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

func uitoa(i uint64) string {
	return fmt.Sprintf("%d", i)
}
