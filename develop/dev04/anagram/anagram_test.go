package anagram

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindAnagrams(t *testing.T) {
	testCases := []struct {
		name     string
		input    []string
		expected map[string][]string
	}{
		{
			name: "basic usage",
			input: []string{"тяпка", "пятак", "пятка", "тяпкА",
				"листок", "слиток", "столик", "ЛИСток", "Го"},
			expected: map[string][]string{
				"листок": {"листок", "слиток", "столик"},
				"тяпка":  {"пятак", "пятка", "тяпка"},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res := FindAnagrams(tc.input)
			assert.Equal(t, len(tc.expected), len(res))
			for k, v := range tc.expected {
				value, has := res[k]
				assert.Greater(t, len(value), 1)
				assert.True(t, has)
				assert.Equal(t, v, value)

				n := len(value)
				for i := 1; i < n; i++ {
					assert.GreaterOrEqual(t, value[i], value[i-1])
				}
			}
		})
	}
}
