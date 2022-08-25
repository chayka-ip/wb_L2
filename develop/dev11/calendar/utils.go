package calendar

import "time"

//DateFromString extracts date from string
func DateFromString(s string) (time.Time, error) {
	date, err := time.Parse(TimeLayout, s)
	if err != nil {
		return time.Time{}, err
	}
	return date, nil
}

//DateFromYearMonthDay ...
func DateFromYearMonthDay(year int, month time.Month, day int) time.Time {
	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
}

//DateFromTime ...
func DateFromTime(t time.Time) time.Time {
	return DateFromYearMonthDay(t.Date())
}

//DateNow trunctates time from time.Time to get clear data
func DateNow() time.Time {
	return DateFromTime(time.Now())
}

//AreEqualEventSlicesByDateAndCount checks whether two event slices are same length
//and have corresponding event with same date in other slice
func AreEqualEventSlicesByDateAndCount(a, b []Event) bool {
	if len(a) != len(b) {
		return false
	}
	foundInd := make(map[int]struct{}, len(b))

	for _, vA := range a {
		for i, vB := range b {
			if _, has := foundInd[i]; !has {
				if vA.Date == vB.Date {
					foundInd[i] = struct{}{}
					break
				}
			}
		}
	}
	return len(a) == len(foundInd)
}
