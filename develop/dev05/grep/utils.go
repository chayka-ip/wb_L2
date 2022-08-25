package grep

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"
)

func maxInt(a, b int) int {
	if a >= b {
		return a
	}
	return b
}

func minInt(a, b int) int {
	if a <= b {
		return a
	}
	return b
}

func readData(readCloser io.ReadCloser) ([]string, error) {
	defer readCloser.Close()
	scanner := bufio.NewScanner(readCloser)

	out := make([]string, 0)

	for scanner.Scan() {
		out = append(out, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return out, nil
}

func printData(src []string, indices []int, printLineNum bool) {
	last := len(src) - 1
	for _, index := range indices {
		if index > last || index < 0 {
			panic("grep: line indices don't match with the data")
		}
		line := src[index]

		if printLineNum {
			fmt.Printf("%d:%s\n", index+1, line)
			continue
		}
		fmt.Println(line)
	}
}

func printCount(n int) {
	fmt.Println(n)
}

func getPattern(rawPattern string, isExact, ignoreCase bool) (string, error) {
	if !isExact {
		return getRegexpPattern(rawPattern, ignoreCase)
	}
	return getExactPattern(rawPattern, ignoreCase), nil
}

func getExactPattern(rawPattern string, ignoreCase bool) string {
	if ignoreCase {
		return strings.ToLower(rawPattern)
	}
	return rawPattern
}

func getRegexpPattern(rawPattern string, ignoreCase bool) (string, error) {
	f := "(?i)"

	if ignoreCase && !strings.HasPrefix(rawPattern, f) {
		rawPattern = f + rawPattern
	}

	if _, err := regexp.Compile(rawPattern); err != nil {
		return "", nil
	}
	return rawPattern, nil
}
