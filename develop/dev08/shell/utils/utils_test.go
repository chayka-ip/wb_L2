package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetFirstArgAndRestAsStr(t *testing.T) {

	testCases := []struct {
		name        string
		input       string
		expected    [2]string
		expectedErr error
	}{
		{
			name:     "basic",
			input:    "  test   other args",
			expected: [2]string{"test", "   other args"},
		},
	}

	for _, tc := range testCases {
		t.Run(t.Name(), func(t *testing.T) {
			cmdName, args, err := GetFirstArgAndRestAsStr(tc.input)
			if tc.expectedErr != nil {
				assert.Equal(t, tc.expectedErr, err)
				return
			}
			assert.Equal(t, tc.expected, [2]string{cmdName, args})
		})
	}

}
