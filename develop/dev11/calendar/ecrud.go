package calendar

//EventCRUD ...
type EventCRUD struct {
	UserID      *uint64 `json:"user_id,omitempty"`
	UID         *uint64 `json:"uid,omitempty"`
	Date        *string `json:"date,omitempty"`
	Title       *string `json:"title,omitempty"`
	Description *string `json:"description,omitempty"`
}

//HasUserID ...
func (e *EventCRUD) HasUserID() bool {
	return e.UserID != nil
}

//HasUID ...
func (e *EventCRUD) HasUID() bool {
	return e.UID != nil
}

//HasDate ...
func (e *EventCRUD) HasDate() bool {
	return e.Date != nil
}

//HasTitle ...
func (e *EventCRUD) HasTitle() bool {
	return e.Title != nil
}

//HasDescription ...
func (e *EventCRUD) HasDescription() bool {
	return e.Description != nil
}

//ToCalendarEvent ...
func (e *EventCRUD) ToCalendarEvent() (*Event, error) {
	out := &Event{}

	if e.UID != nil {
		out.UID = *e.UID
	}
	if e.Title != nil {
		out.Title = *e.Title
	}
	if e.Description != nil {
		out.Description = *e.Description
	}

	if e.Date != nil {
		d, err := DateFromString(*e.Date)
		if err != nil {
			return nil, err
		}
		out.Date = d
	}
	return out, nil
}
