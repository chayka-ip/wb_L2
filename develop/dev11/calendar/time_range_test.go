package calendar

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	dayWithoutOneNanoSec = 24*time.Hour - 1
)

func TestIsTimeInRange(t *testing.T) {
	type testCaseData struct {
		name      string
		time      time.Time
		timeRange TimeRange
		isInRange bool
	}
	refTimeRange := TimeRange{
		MinTime: DateFromYearMonthDay(2020, 12, 1).Add(dayWithoutOneNanoSec),
		MaxTime: DateFromYearMonthDay(2020, 12, 3),
	}
	testCases := []testCaseData{
		{
			name:      "in range",
			time:      time.Date(2020, 12, 2, 1, 1, 1, 1, time.UTC),
			timeRange: refTimeRange,
			isInRange: true,
		},
		{
			name:      "not in range",
			time:      DateFromYearMonthDay(2020, 12, 1).Add(dayWithoutOneNanoSec),
			timeRange: refTimeRange,
			isInRange: false,
		},
		{
			name:      "not in range",
			time:      DateFromYearMonthDay(2020, 12, 3),
			timeRange: refTimeRange,
			isInRange: false,
		},
		{
			name:      "not in range",
			time:      DateFromYearMonthDay(2020, 12, 4),
			timeRange: refTimeRange,
			isInRange: false,
		},
		{
			name:      "not in range",
			time:      DateFromYearMonthDay(2020, 11, 4),
			timeRange: refTimeRange,
			isInRange: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res := tc.timeRange.IsTimeInInterval(tc.time)
			assert.Equal(t, tc.isInRange, res)
		})
	}
}

func TestMakeTimeRangeForDay(t *testing.T) {
	tt := time.Date(2020, 12, 2, 1, 1, 1, 1, time.UTC)
	testCases := []struct {
		name       string
		time       time.Time
		want       TimeRange
		expectedEq bool
	}{
		{
			name: "valid result",
			time: tt,
			want: TimeRange{
				MinTime: DateFromYearMonthDay(2020, 12, 1).Add(dayWithoutOneNanoSec),
				MaxTime: DateFromYearMonthDay(2020, 12, 3),
			},
			expectedEq: true,
		},
		{
			name: "wrong result",
			time: tt,
			want: TimeRange{
				MinTime: DateFromYearMonthDay(2020, 12, 1).Add(dayWithoutOneNanoSec - 1),
				MaxTime: DateFromYearMonthDay(2020, 12, 3),
			},
			expectedEq: false,
		},
		{
			name: "wrong result 2",
			time: tt,
			want: TimeRange{
				MinTime: DateFromYearMonthDay(2020, 12, 1).Add(dayWithoutOneNanoSec),
				MaxTime: DateFromYearMonthDay(2020, 12, 3).Add(1),
			},
			expectedEq: false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			get := makeTimeRangeForDay(tc.time)
			eq := reflect.DeepEqual(tc.want, get)
			assert.Equal(t, tc.expectedEq, eq)
		})
	}
}

func TestMakeTimeRangeForWeek(t *testing.T) {
	testCases := []struct {
		name       string
		time       time.Time
		want       TimeRange
		expectedEq bool
	}{
		{
			name: "valid result",
			time: time.Date(2020, 12, 12, 1, 1, 1, 1, time.UTC),
			want: TimeRange{
				MinTime: DateFromYearMonthDay(2020, 12, 6).Add(dayWithoutOneNanoSec),
				MaxTime: DateFromYearMonthDay(2020, 12, 14),
			},
			expectedEq: true,
		},
		{
			name: "valid result",
			time: time.Date(2020, 12, 13, 1, 1, 1, 1, time.UTC),
			want: TimeRange{
				MinTime: DateFromYearMonthDay(2020, 12, 6).Add(dayWithoutOneNanoSec),
				MaxTime: DateFromYearMonthDay(2020, 12, 14),
			},
			expectedEq: true,
		},
		{
			name: "valid result",
			time: time.Date(2020, 12, 10, 1, 1, 1, 1, time.UTC),
			want: TimeRange{
				MinTime: DateFromYearMonthDay(2020, 12, 6).Add(dayWithoutOneNanoSec),
				MaxTime: DateFromYearMonthDay(2020, 12, 14),
			},
			expectedEq: true,
		},
		{
			name: "wrong result",
			time: time.Date(2020, 12, 10, 1, 1, 1, 1, time.UTC),
			want: TimeRange{
				MinTime: DateFromYearMonthDay(2020, 12, 6).Add(dayWithoutOneNanoSec + 1),
				MaxTime: DateFromYearMonthDay(2020, 12, 14),
			},
			expectedEq: false,
		},
		{
			name: "wrong result",
			time: time.Date(2020, 12, 10, 1, 1, 1, 1, time.UTC),
			want: TimeRange{
				MinTime: DateFromYearMonthDay(2020, 12, 6).Add(dayWithoutOneNanoSec),
				MaxTime: DateFromYearMonthDay(2020, 12, 14).Add(-1),
			},
			expectedEq: false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			get := makeTimeRangeForWeek(tc.time)
			eq := reflect.DeepEqual(tc.want, get)
			fmt.Println(get)
			assert.Equal(t, tc.expectedEq, eq)
		})
	}
}

func TestMakeTimeRangeForMonth(t *testing.T) {
	tt := time.Date(2020, 12, 2, 1, 1, 1, 1, time.UTC)
	testCases := []struct {
		name       string
		time       time.Time
		want       TimeRange
		expectedEq bool
	}{
		{
			name: "valid result",
			time: tt,
			want: TimeRange{
				MinTime: DateFromYearMonthDay(2020, 11, 30).Add(dayWithoutOneNanoSec),
				MaxTime: DateFromYearMonthDay(2021, 1, 1),
			},
			expectedEq: true,
		},
		{
			name: "wrong result",
			time: time.Date(2020, 12, 2, 1, 1, 1, 1, time.UTC),
			want: TimeRange{
				MinTime: DateFromYearMonthDay(2020, 11, 30).Add(dayWithoutOneNanoSec + 1),
				MaxTime: DateFromYearMonthDay(2021, 1, 1),
			},
			expectedEq: false,
		},
		{
			name: "wrong result 2",
			time: time.Date(2020, 12, 2, 1, 1, 1, 1, time.UTC),
			want: TimeRange{
				MinTime: DateFromYearMonthDay(2020, 11, 30),
				MaxTime: DateFromYearMonthDay(2021, 1, 1).Add(-1),
			},
			expectedEq: false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			get := makeTimeRangeForMonth(tc.time)
			eq := reflect.DeepEqual(tc.want, get)
			assert.Equal(t, tc.expectedEq, eq)
		})
	}
}
