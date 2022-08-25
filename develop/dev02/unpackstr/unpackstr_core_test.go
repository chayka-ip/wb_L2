package unpackstr

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnpackString(t *testing.T) {
	testCases := []struct {
		name        string
		src         string
		expected    string
		expectedErr error
	}{
		{
			name:        "basic",
			src:         "a4bc2d5e",
			expected:    "aaaabccddddde",
			expectedErr: nil,
		},
		{
			name:        "no repeats",
			src:         "abcd",
			expected:    "abcd",
			expectedErr: nil,
		},
		{
			name:        "starts with digit",
			src:         "45",
			expected:    "",
			expectedErr: ErrInvalidString,
		},
		{
			name:        "empty",
			src:         "",
			expected:    "",
			expectedErr: nil,
		},
		{
			name:        "escape escape",
			src:         `\\`,
			expected:    `\`,
			expectedErr: nil,
		},
		{
			name:        "unpacking escapes",
			src:         `qwe\\5`,
			expected:    `qwe\\\\\`,
			expectedErr: nil,
		},
		{
			name:        "empty escape is not allowed",
			src:         `a\`,
			expected:    "",
			expectedErr: ErrInvalidString,
		},
		{
			name:        "letter escape is not allowed",
			src:         `a\a`,
			expected:    "",
			expectedErr: ErrInvalidString,
		},
		{
			name:        "escape numbers",
			src:         `qwe\4\5`,
			expected:    `qwe45`,
			expectedErr: nil,
		},
		{
			name:        "unpacking numbers",
			src:         `qwe\45`,
			expected:    `qwe44444`,
			expectedErr: nil,
		},
		{
			name:        "space characters are not allowed",
			src:         `q w`,
			expected:    ``,
			expectedErr: ErrInvalidString,
		},
		{
			name:        "single space is not allowed",
			src:         ` `,
			expected:    ``,
			expectedErr: ErrInvalidString,
		},
		{
			name:        "single escape is not allowed",
			src:         `\`,
			expected:    ``,
			expectedErr: ErrInvalidString,
		},
		{
			name:        "digit with last escape",
			src:         `t4\5`,
			expected:    `tttt5`,
			expectedErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := unpack(tc.src)
			if err != tc.expectedErr {
				log.Fatal(err)
			}
			assert.Equal(t, tc.expected, res)
		})
	}

}
