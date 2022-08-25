package utils

import (
	"errors"
	"os"
	"strings"
	"unicode"
)

var (
	//ErrBadData is returned when provided data could not be handled correctly
	ErrBadData = errors.New("bad data")
)

const (
	PIPE = "|"
)

//GetFirstArgAndRestAsStr returns first argument and other args as string
func GetFirstArgAndRestAsStr(s string) (string, string, error) {
	s = strings.TrimLeft(s, " ")

	var otherArgs string
	var cmdName strings.Builder
	cmdName.Grow(10)

	if len(s) > 0 {
		for byteInd, r := range s {
			if unicode.IsSpace(r) {
				otherArgs = s[byteInd:]
				break
			}
			cmdName.WriteRune(r)
		}

		return cmdName.String(), otherArgs, nil
	}
	return "", "", ErrBadData
}

//ValidateNumArgs ...
func ValidateNumArgs(min, max int, args string, lessArgErr, greaterArgsErr error) error {
	n := len(strings.Fields(args))
	if n < min {
		return lessArgErr
	}
	if n > max {
		return greaterArgsErr
	}
	return nil
}

//HasPipe ...
func HasPipe(s string) bool {
	return strings.Contains(s, PIPE)
}

//TryOpenDir ...
func TryOpenDir(path string) error {
	_, err := os.ReadDir(path)
	return err
}
