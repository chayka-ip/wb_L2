package sort

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type colSortTestData struct {
	name     string
	data     []string
	expected []string
	opt      *options
}

func TestDoStingSort(t *testing.T) {
	baseInputData := []string{
		"The second Klan started in 1915",
		"as a small group in Georgia",
		"Calvin Jones, and James Crowe",
		"please, stop it",
		"The word had previously been",
		"zzz fun a zz",
		"get some help",
		"other fraternal organizations",
		"and sometimes still recognized",
		"the attackers by voice and mannerisms",
	}

	testCases := []struct {
		name         string
		data         []string
		expected     []string
		isDescending bool
	}{
		{
			name: "default sort strings asc",
			data: baseInputData,
			expected: []string{
				"Calvin Jones, and James Crowe",
				"The second Klan started in 1915",
				"The word had previously been",
				"and sometimes still recognized",
				"as a small group in Georgia",
				"get some help",
				"other fraternal organizations",
				"please, stop it",
				"the attackers by voice and mannerisms",
				"zzz fun a zz",
			},
			isDescending: false,
		},
		{
			name: "default sort strings desc",
			data: baseInputData,
			expected: []string{
				"zzz fun a zz",
				"the attackers by voice and mannerisms",
				"please, stop it",
				"other fraternal organizations",
				"get some help",
				"as a small group in Georgia",
				"and sometimes still recognized",
				"The word had previously been",
				"The second Klan started in 1915",
				"Calvin Jones, and James Crowe",
			},
			isDescending: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			d := append([]string{}, tc.data...)
			res := doDefaultStringSort(d, tc.isDescending)
			assert.Equal(t, tc.expected, res)
		})
	}
}

func TestDoNumericStingSort(t *testing.T) {
	baseInputData := []string{
		"dddd",
		"1ff",
		"a2",
		"a111",
		"11ddd",
		"sss",
		"2f",
	}

	testCases := []struct {
		name         string
		data         []string
		expected     []string
		isDescending bool
	}{
		{
			name: "asc",
			data: baseInputData,
			expected: []string{
				"a111",
				"a2",
				"dddd",
				"sss",
				"1ff",
				"2f",
				"11ddd",
			},
			isDescending: false,
		},
		{
			name: "desc",
			data: baseInputData,
			expected: []string{
				"11ddd",
				"2f",
				"1ff",
				"sss",
				"dddd",
				"a2",
				"a111",
			},
			isDescending: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			d := append([]string{}, tc.data...)
			res := doStringNumericSort(d, tc.isDescending)
			assert.Equal(t, tc.expected, res)
		})
	}
}

func TestColumnSortString(t *testing.T) {
	baseInputData := []string{
		"fsdsda opsd aso asdpo",
		"cspds gpfdp psdfd",
		"zals",
		"bbsdfpd ukdsk",
		"fsdsda zpsd naso bcdpo",
		"fsdsda gpsd khso fdpo",
	}
	testCases := []colSortTestData{
		{
			name: "sort 4 col asc",
			data: baseInputData,
			expected: []string{
				"bbsdfpd ukdsk",
				"cspds gpfdp psdfd",
				"zals",
				"fsdsda opsd aso asdpo",
				"fsdsda zpsd naso bcdpo",
				"fsdsda gpsd khso fdpo",
			},
			opt: &options{
				isDescOrder:  false,
				sortColIndex: 3,
				separator:    defaultSeparator,
			},
		},
		{
			name: "sort 4 col desc",
			data: baseInputData,
			expected: []string{
				"fsdsda gpsd khso fdpo",
				"fsdsda zpsd naso bcdpo",
				"fsdsda opsd aso asdpo",
				"zals",
				"cspds gpfdp psdfd",
				"bbsdfpd ukdsk",
			},
			opt: &options{
				isDescOrder:  true,
				sortColIndex: 3,
				separator:    defaultSeparator,
			},
		},
		{
			name: "sort 3 col asc",
			data: baseInputData,
			expected: []string{
				"bbsdfpd ukdsk",
				"zals",
				"fsdsda opsd aso asdpo",
				"fsdsda gpsd khso fdpo",
				"fsdsda zpsd naso bcdpo",
				"cspds gpfdp psdfd",
			},
			opt: &options{
				isDescOrder:  false,
				sortColIndex: 2,
				separator:    defaultSeparator,
			},
		},
		{
			name: "sort 3 col desc",
			data: baseInputData,
			expected: []string{
				"cspds gpfdp psdfd",
				"fsdsda zpsd naso bcdpo",
				"fsdsda gpsd khso fdpo",
				"fsdsda opsd aso asdpo",
				"zals",
				"bbsdfpd ukdsk",
			},
			opt: &options{
				isDescOrder:  true,
				sortColIndex: 2,
				separator:    defaultSeparator,
			},
		},
		{
			name: "sort 2 col asc",
			data: baseInputData,
			expected: []string{
				"zals",
				"cspds gpfdp psdfd",
				"fsdsda gpsd khso fdpo",
				"fsdsda opsd aso asdpo",
				"bbsdfpd ukdsk",
				"fsdsda zpsd naso bcdpo",
			},
			opt: &options{
				isDescOrder:  false,
				sortColIndex: 1,
				separator:    defaultSeparator,
			},
		},
		{
			name: "sort 2 col desc",
			data: baseInputData,
			expected: []string{
				"fsdsda zpsd naso bcdpo",
				"bbsdfpd ukdsk",
				"fsdsda opsd aso asdpo",
				"fsdsda gpsd khso fdpo",
				"cspds gpfdp psdfd",
				"zals",
			},
			opt: &options{
				isDescOrder:  true,
				sortColIndex: 1,
				separator:    defaultSeparator,
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			d := append([]string{}, tc.data...)
			res := doColumnSort(tc.opt, d)
			assert.Equal(t, tc.expected, res)
		})
	}
}

func TestColumnSortNumeric(t *testing.T) {
	baseInputData := []string{
		"sdb 22 04",
		"4 85",
		"k 15 11",
		"1 13",
		"9",
	}
	testCases := []colSortTestData{
		{
			name: "sort 3 col asc",
			data: baseInputData,
			expected: []string{
				"1 13",
				"4 85",
				"9",
				"sdb 22 04",
				"k 15 11",
			},
			opt: &options{
				isDescOrder:  false,
				sortColIndex: 2,
				sortType:     sortTypeNumeric,
				separator:    defaultSeparator,
			},
		},
		{
			name: "sort 3 col desc",
			data: baseInputData,
			expected: []string{
				"k 15 11",
				"sdb 22 04",
				"9",
				"4 85",
				"1 13",
			},
			opt: &options{
				isDescOrder:  true,
				sortColIndex: 2,
				sortType:     sortTypeNumeric,
				separator:    defaultSeparator,
			},
		},
		{
			name: "sort 2 col asc",
			data: baseInputData,
			expected: []string{
				"9",
				"1 13",
				"k 15 11",
				"sdb 22 04",
				"4 85",
			},
			opt: &options{
				isDescOrder:  false,
				sortColIndex: 1,
				sortType:     sortTypeNumeric,
				separator:    defaultSeparator,
			},
		},
		{
			name: "sort 2 col desc",
			data: baseInputData,
			expected: []string{
				"4 85",
				"sdb 22 04",
				"k 15 11",
				"1 13",
				"9",
			},
			opt: &options{
				isDescOrder:  true,
				sortColIndex: 1,
				sortType:     sortTypeNumeric,
				separator:    defaultSeparator,
			},
		},
		{
			name: "sort 1 col asc",
			data: baseInputData,
			expected: []string{
				"k 15 11",
				"sdb 22 04",
				"1 13",
				"4 85",
				"9",
			},
			opt: &options{
				isDescOrder:  false,
				sortColIndex: 0,
				sortType:     sortTypeNumeric,
				separator:    defaultSeparator,
			},
		},
		{
			name: "sort 1 col desc",
			data: baseInputData,
			expected: []string{
				"9",
				"4 85",
				"1 13",
				"sdb 22 04",
				"k 15 11",
			},
			opt: &options{
				isDescOrder:  true,
				sortColIndex: 0,
				sortType:     sortTypeNumeric,
				separator:    defaultSeparator,
			},
		},
		{
			name: "correct middle split",
			data: []string{
				"1 02 3",
				"1 03             4",
				"1 001 5",
			},
			expected: []string{
				"1 02 3",
				"1 03             4",
				"1 001 5",
			},
			opt: &options{
				isDescOrder:  false,
				sortColIndex: 2,
				sortType:     sortTypeNumeric,
				separator:    defaultSeparator,
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			d := append([]string{}, tc.data...)
			res := doColumnSort(tc.opt, d)
			assert.Equal(t, tc.expected, res)
		})
	}
}

func TestRemoveDuplicates(t *testing.T) {
	data := []string{
		"a",
		"a",
		"b",
	}
	expected := []string{
		"a",
		"b",
	}
	t.Run("remove duplicates", func(t *testing.T) {
		res := removeDuplicates(data)
		assert.Equal(t, expected, res)
	})
}

func TestSplitStringAndCleanUp(t *testing.T) {
	sep := " "
	src := " 1 03             4 "
	expected := []string{" 1", "03", "4 "}

	t.Run("separation and cleanup", func(t *testing.T) {
		res := splitStringAndCleanUp(src, sep)
		assert.Equal(t, expected, res)
	})
}
