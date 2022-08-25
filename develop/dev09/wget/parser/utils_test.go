package parser

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRuneSliceEndsWith(t *testing.T) {
	testCases := []struct {
		name     string
		src      string
		end      string
		expected bool
	}{
		{
			name:     "ok 1",
			src:      "abbbddd",
			end:      "ddd",
			expected: true,
		},
		{
			name:     "ok 2",
			src:      "abbbddd",
			end:      "",
			expected: true,
		},
		{
			name:     "nope 1",
			src:      "abbbddd",
			end:      "b",
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := runeSliceEndsWith([]rune(tc.src), []rune(tc.end))
			assert.Equal(t, tc.expected, r)
		})
	}
}

func TestGetNumLevelsToCommonParentDir(t *testing.T) {
	testCases := []struct {
		name        string
		pathRef     string
		pathTarget  string
		expectedNum int
		expectedErr bool
	}{
		{
			name:        "basic same dir",
			pathRef:     "a/b/c/d",
			pathTarget:  "a/b/c/d",
			expectedNum: 0,
		},
		{
			name:        "basic path 1 longer",
			pathRef:     "a/b/c/d",
			pathTarget:  "a/b",
			expectedNum: 2,
		},
		{
			name:        "basic same dir",
			pathRef:     "a/b/c/d",
			pathTarget:  "a/b/e/f",
			expectedNum: 2,
		},
		{
			name:        "basic path 2 longer",
			pathRef:     "a/b",
			pathTarget:  "a/b/c/d",
			expectedNum: 0,
		},
		{
			name:        "basic path 2 longer",
			pathRef:     "a/b",
			pathTarget:  "a/f/c/d",
			expectedNum: 1,
		},
		{
			name:        "basic path 2 longer",
			pathRef:     "a/b",
			pathTarget:  "n/f/c/d",
			expectedErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			n, err := getNumLevelsToCommonParentDir(tc.pathRef, tc.pathTarget)
			if tc.expectedErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedNum, n)
		})
	}
}

func TestResourceURLToLocalRelative(t *testing.T) {
	testCases := []struct {
		name         string
		pathTarget   string
		pathResource string
		expected     string
		expectedErr  bool
	}{
		{
			name:         "basic",
			pathTarget:   "a/b/c/d",
			pathResource: "a/b/g/f",
			expected:     "../../g/f",
		},
		{
			name:         "zero levels",
			pathTarget:   "a/b",
			pathResource: "a/b/f/g",
			expected:     "/f/g",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := resourceLinkToLocalRelative(tc.pathTarget, tc.pathResource)
			if tc.expectedErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, res)
		})
	}
}

func TestOverrideBytes(t *testing.T) {
	testCases := []struct {
		name        string
		src         string
		data        string
		startPos    int
		byteDelta   int
		expected    string
		expectedErr bool
	}{
		{
			name:      "basic no delta",
			src:       strings.Repeat("a", 10),
			data:      "bbb",
			startPos:  3,
			byteDelta: 0,
			expected:  "aaabbbaaaa",
		},
		{
			name:      "inc delta",
			src:       strings.Repeat("a", 10),
			data:      "bbb",
			startPos:  3,
			byteDelta: 3,
			expected:  "aaabbbaaaaaaa",
		},
		{
			name:      "dec delta",
			src:       strings.Repeat("a", 10),
			data:      "bbb",
			startPos:  3,
			byteDelta: -3,
			expected:  "aaabbba",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := overrideBytes([]byte(tc.src), []byte(tc.data), tc.startPos, tc.byteDelta)
			if tc.expectedErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, string(res))
		})
	}
}
