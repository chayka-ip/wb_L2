package cut

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseFields(t *testing.T) {
	testCases := []struct {
		name        string
		input       string
		expected    [][2]int
		expectedErr error
	}{
		{
			name:     "comma sep",
			input:    " 1,2, 5  ",
			expected: [][2]int{{1, 1}, {2, 2}, {5, 5}},
		},
		{
			name:     "hyphen sep",
			input:    " 5 -10 ",
			expected: [][2]int{{5, 10}},
		},
		{
			name:     "basic",
			input:    " 1,5-7",
			expected: [][2]int{{1, 1}, {5, 7}},
		},
		{
			name:     "with lower ommit",
			input:    " 1,-7",
			expected: [][2]int{{1, 1}, {minSearchRange, 7}},
		},
		{
			name:     "with upper ommit",
			input:    " 1,1-",
			expected: [][2]int{{1, 1}, {1, maxSearchRange}},
		},
		{
			name:     "empty is max range",
			input:    " ",
			expected: [][2]int{{minSearchRange, maxSearchRange}},
		},
		{
			name:        "only comma is error",
			input:       " ,",
			expected:    nil,
			expectedErr: ErrBadFields,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := parseFields(tc.input)
			if tc.expectedErr != nil {
				assert.Equal(t, tc.expectedErr, err)
				return
			}
			assert.NoError(t, err)
			for _, v := range tc.expected {
				has := res.hasRange(v[0], v[1])
				assert.True(t, has)
			}
		})
	}
}

func TestUnpackIntRangeEnds(t *testing.T) {
	testCases := []struct {
		name        string
		input       string
		expectedA   int
		expectedB   int
		expectedErr error
	}{
		{
			name:        "basic",
			input:       "1-10",
			expectedA:   1,
			expectedB:   10,
			expectedErr: nil,
		},
		{
			name:        "ommit a",
			input:       "-10",
			expectedA:   minSearchRange,
			expectedB:   10,
			expectedErr: nil,
		},
		{
			name:        "ommit b",
			input:       "5-",
			expectedA:   5,
			expectedB:   maxSearchRange,
			expectedErr: nil,
		},
		{
			name:        "ommit all gives err",
			input:       "-",
			expectedA:   0,
			expectedB:   0,
			expectedErr: ErrBadFields,
		},
		{
			name:        "a > b is error",
			input:       "10-1",
			expectedA:   0,
			expectedB:   0,
			expectedErr: ErrBadFields,
		},
		{
			name:        "empty is error",
			input:       "",
			expectedA:   0,
			expectedB:   0,
			expectedErr: ErrBadFields,
		},
		{
			name:        "a = b is ok",
			input:       "10-10",
			expectedA:   10,
			expectedB:   10,
			expectedErr: nil,
		},
		{
			name:        "a less than min is error",
			input:       strconv.Itoa(minSearchRange-1) + "-" + strconv.Itoa(minSearchRange+1),
			expectedA:   0,
			expectedB:   0,
			expectedErr: ErrBadFields,
		},
		{
			name:        "b greater than max is error",
			input:       strconv.Itoa(minSearchRange) + "-" + strconv.Itoa(maxSearchRange+1),
			expectedA:   0,
			expectedB:   0,
			expectedErr: ErrBadFields,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			a, b, err := unpackIntRangeEnds(tc.input)
			if tc.expectedErr != nil {
				assert.Equal(t, tc.expectedErr, err)
				return
			}
			assert.Equal(t, tc.expectedErr, err)
			assert.Equal(t, tc.expectedA, a)
			assert.Equal(t, tc.expectedB, b)
		})
	}
}
