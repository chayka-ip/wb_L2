package sort

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"unicode"
)

func splitStringAndCleanUp(s, sep string) []string {
	sp := strings.Split(s, sep)
	n := len(sp)
	if n == 1 {
		return sp
	}
	last := n - 1
	firtsIsSep := false
	out := []string{}
	for i, v := range sp {
		if i == 0 {
			if v == "" {
				firtsIsSep = true
				continue
			}
			out = append(out, v)
			continue
		}
		if i == 1 {
			if firtsIsSep {
				out = append(out, joinSep(v, sep, true))
				continue
			}
		}
		if i == last {
			if v == "" {
				outLast := len(out) - 1
				out[outLast] = joinSep(out[outLast], sep, false)
				continue
			}
			out = append(out, v)
			continue
		}

		if v != "" {
			out = append(out, v)
		}
	}
	return out
}

func joinSep(s, sep string, sepInfront bool) string {
	if sepInfront {
		return sep + s
	}
	return s + sep
}

func startsWithDigit(s string) bool {
	if s == "" {
		return false
	}
	return unicode.IsDigit(rune(s[0]))
}

func getNumberFromStringStart(s string) (int, error) {
	badret := func() (int, error) { return 0, errStrNotStartsWithNumber }
	if s == "" {
		return badret()
	}
	var b strings.Builder

	hasStartZero, hasNonZero := false, false

	for _, r := range []rune(s) {
		if !unicode.IsDigit(r) {
			break
		}
		if r == '0' && !hasNonZero {
			hasStartZero = true
			continue
		}
		hasNonZero = true
		b.WriteRune(r)
	}
	n, err := strconv.Atoi(b.String())
	if err != nil {
		if hasStartZero {
			return 0, nil
		}
		return badret()
	}
	return n, nil
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

func writeData(data []string, filePath string) error {
	if filePath == "" {
		writeToStdOut(data)
		return nil
	}

	file, err := os.Create(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to create file %s: %v\n", filePath, err)
		return errBadOpenFile
	}
	for i, v := range data {
		s := v
		if i > 0 {
			s = fmt.Sprintf("\n%s", v)
		}
		_, err := file.WriteString(s)
		if err != nil {
			fmt.Fprintf(os.Stderr, "unable to write to file file %s: %v\n", filePath, err)
			return errWriteFail
		}
	}
	return nil
}

func writeToStdOut(data []string) {
	for _, v := range data {
		fmt.Fprintf(os.Stdout, "%s\n", v)
	}
}

func removeDuplicates(data []string) []string {
	var exists = make(map[string]struct{})
	var out = make([]string, 0)

	for _, v := range data {
		if _, has := exists[v]; !has {
			exists[v] = struct{}{}
			out = append(out, v)
		}
	}

	return out
}
