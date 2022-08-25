package grep

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsStringMatches(t *testing.T) {

	testCases := []struct {
		name     string
		data     []string
		opt      *options
		expected bool
	}{
		{
			name: "exact not ignore case (1)",
			data: []string{"aa bbb dd CCC"},
			opt: &options{
				pattern:      "dd",
				isExactMatch: true,
				isIgnoreCase: false,
			},
			expected: true,
		},
		{
			name: "exact not ignore case (2)",
			data: []string{"aa bbb dd CCC"},
			opt: &options{
				pattern:      "cc",
				isExactMatch: true,
				isIgnoreCase: false,
			},
			expected: false,
		},
		{
			name: "exact ignore case (1)",
			data: []string{"aa bbb dd CCC"},
			opt: &options{
				pattern:      "cc",
				isExactMatch: true,
				isIgnoreCase: true,
			},
			expected: true,
		},
		{
			name: "regexp not ignore case (1)",
			data: []string{"aa bbb dd CCC"},
			opt: &options{
				pattern:      "[b-b]",
				isExactMatch: false,
				isIgnoreCase: false,
			},
			expected: true,
		},
		{
			name: "regexp not ignore case (2)",
			data: []string{"aa bbb dd CCC"},
			opt: &options{
				pattern:      "[B-B]",
				isExactMatch: false,
				isIgnoreCase: false,
			},
			expected: false,
		},
		{
			name: "regexp ignore case (1)",
			data: []string{"aa bbb dd CCC"},
			opt: &options{
				pattern:      "[B-B]",
				isExactMatch: false,
				isIgnoreCase: true,
			},
			expected: true,
		},
		{
			name: "regexp ivnert search (1)",
			data: []string{"aa bbb dd CCC"},
			opt: &options{
				pattern:        "[B-B]",
				isExactMatch:   false,
				isIgnoreCase:   false,
				isInvertSearch: true,
			},
			expected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res := findMatchedLineIndices(tc.opt, tc.data)
			assert.Equal(t, tc.expected, len(res) == len(tc.data))
		})
	}
}

func TestGetLineIndicesToPrint(t *testing.T) {

	testCases := []struct {
		name     string
		data     []int
		before   int
		after    int
		total    int
		expected []int
	}{
		{
			name:     "basic",
			data:     []int{5, 10},
			before:   2,
			after:    2,
			total:    100,
			expected: []int{3, 4, 5, 6, 7, 8, 9, 10, 11, 12},
		},
		{
			name:     "basic 2",
			data:     []int{5, 10},
			before:   1,
			after:    1,
			total:    100,
			expected: []int{4, 5, 6, 9, 10, 11},
		},
		{
			name:     "full range",
			data:     []int{5, 10},
			before:   10,
			after:    10,
			total:    10,
			expected: []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
		},
		{
			name:     "current only",
			data:     []int{10},
			before:   0,
			after:    0,
			total:    100,
			expected: []int{10},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res := getLineIndicesToPrint(tc.before, tc.after, tc.total, tc.data)
			assert.Equal(t, tc.expected, res)
		})
	}
}
