package grep

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"
)

//ExecuteCLI executes grep command with arguments
func ExecuteCLI(args []string) int {
	opt, err := newOptions(args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 2
	}

	if err := run(opt); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}

	return 0
}

func run(opt *options) error {
	data, err := readData(opt.reader)
	if err != nil {
		return err
	}

	matches := findMatchedLineIndices(opt, data)

	if opt.isPrintMatchCount {
		printCount(len(matches))
		return nil
	}

	before := opt.numLinesBeforeMatch
	after := opt.numLinesAfterMatch
	total := len(data)
	print := getLineIndicesToPrint(before, after, total, matches)

	printData(data, print, opt.isPrintLineNum)

	return nil
}

func getLineIndicesToPrint(before, after, total int, matches []int) []int {
	k := len(matches) / 3
	m := make(map[int]struct{})
	out := make([]int, 0, k)

	for _, n := range matches {
		lower, upper, _ := getIndexRangeBounds(before, after, total, n)
		for i := lower; i <= upper; i++ {
			if _, has := m[i]; !has {
				m[i] = struct{}{}
				out = append(out, i)
			}
		}
	}

	sort.Ints(out)
	return out
}

func getIndexRangeBounds(before, after, total, current int) (int, int, error) {
	if (current > total) || (current < 0) || (before < 0) || (after < 0) || (total < 1) {
		return 0, 0, errors.New("wrong bounds configuration")
	}
	lower := maxInt(current-before, 0)
	upper := minInt(current+after, total-1)
	return lower, upper, nil
}

func findMatchedLineIndices(opt *options, data []string) []int {
	invert := opt.isInvertSearch
	igoreCase := opt.isIgnoreCase
	isExact := opt.isExactMatch
	if isExact {
		return findMatchedLineIndicesExact(data, opt.GetExactPattern(), igoreCase, invert)
	}
	return findMatchedLineIndicesRegexp(data, opt.GetRegExpPattern(), invert)
}

func isMatch(contains, invert bool) bool {
	return (contains && !invert) || (!contains && invert)
}

func findMatchedLineIndicesExact(data []string, pattern string, ignoreCase, invert bool) []int {
	out := make([]int, 0, len(data)/2)
	if ignoreCase {
		pattern = strings.ToLower(pattern)
	}

	for i, s := range data {
		if ignoreCase {
			s = strings.ToLower(s)
		}
		cont := strings.Contains(s, pattern)
		if isMatch(cont, invert) {
			out = append(out, i)
		}
	}
	return out
}

func findMatchedLineIndicesRegexp(data []string, r *regexp.Regexp, invert bool) []int {
	out := make([]int, 0, len(data)/2)
	for i, s := range data {
		cont := r.MatchString(s)
		if isMatch(cont, invert) {
			out = append(out, i)
		}
	}
	return out
}
