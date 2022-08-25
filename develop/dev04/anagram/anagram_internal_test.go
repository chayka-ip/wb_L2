package anagram

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringToSortedRuneSlice(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected []rune
	}{
		{
			name:     "all chars lower case",
			input:    "абвгдеёжзийклмнопрстуфхцчшщъыьэюя",
			expected: []rune("абвгдежзийклмнопрстуфхцчшщъыьэюяё"),
		},
		{
			name:     "some chars with duplicates lower case",
			input:    "ваббба",
			expected: []rune("аабббв"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res := stringToSortedRuneSlice(tc.input)
			assert.Equal(t, tc.expected, res)
		})
	}

}
