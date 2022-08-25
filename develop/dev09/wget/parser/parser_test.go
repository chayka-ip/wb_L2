package parser

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractLinksFromSring(t *testing.T) {
	testCases := []struct {
		name     string
		src      string
		expected []string
	}{
		{
			name:     "basic 1",
			src:      ` href="https://test1.org"`,
			expected: []string{"https://test1.org"},
		},
		{
			name: "basic 1",
			src: ` href="https://test1.org"
				   src="https://test2.org"
				   content="https://test3.org"`,
			expected: []string{"https://test1.org", "https://test2.org", "https://test3.org"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := ExtractLinksFromString(tc.src)
			eq := reflect.DeepEqual(tc.expected, r.GetDataAsStringSlice())
			assert.True(t, eq)
		})
	}
}
