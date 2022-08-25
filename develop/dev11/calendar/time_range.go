package calendar

import "time"

//TimeRange represents time interval
type TimeRange struct {
	MinTime time.Time
	MaxTime time.Time
}

//IsTimeInInterval checks wheter time belongs to time interval
func (r *TimeRange) IsTimeInInterval(t time.Time) bool {
	return t.After(r.MinTime) && t.Before(r.MaxTime)
}

//GetTimeRangeFromDateRange calculates time interval based on DateRange provided.
//Boundaries are not included in the interval
func GetTimeRangeFromDateRange(t time.Time, rangeType EventDateRange) TimeRange {
	switch rangeType {
	case EventDAY:
		return makeTimeRangeForDay(t)
	case EventWEEK:
		return makeTimeRangeForWeek(t)
	case EventMONTH:
		return makeTimeRangeForMonth(t)
	default:
		return TimeRange{}
	}
}

func makeTimeRangeForDay(t time.Time) TimeRange {
	year, month, day := t.Date()
	t = DateFromYearMonthDay(year, month, day)
	min := t.Add(-1)
	max := t.AddDate(0, 0, 1)
	return TimeRange{
		MinTime: min,
		MaxTime: max,
	}
}

func makeTimeRangeForWeek(t time.Time) TimeRange {
	year, month, day := t.Date()
	t = DateFromYearMonthDay(year, month, day)

	sub, add := 0, 0
	weekday := t.UTC().Weekday()
	switch weekday {
	case time.Sunday:
		sub, add = 6, 1
	default:
		w := int(weekday)
		sub, add = w-1, 8-w
	}

	min := t.AddDate(0, 0, -sub).Add(-1)
	max := t.AddDate(0, 0, add)

	return TimeRange{
		MinTime: min,
		MaxTime: max,
	}
}

func makeTimeRangeForMonth(t time.Time) TimeRange {
	year, month, _ := t.Date()
	t = DateFromYearMonthDay(year, month, 1)
	min := t.Add(-1)
	max := t.AddDate(0, 1, 0)
	return TimeRange{
		MinTime: min,
		MaxTime: max,
	}
}
