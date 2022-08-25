package unpackstr

import (
	"errors"
	"fmt"
	"os"
	"unpack-str/utils"
)

var (
	//ErrInvalidString is returned when source string is invalid
	ErrInvalidString = errors.New("string is invalid")
)

//unpack unpacks string according to the specification
func unpack(s string) (string, error) {
	if s == "" {
		return passedStringReturn(s)
	}

	runeSlice := []rune(s)
	if utils.IsDigitRune(runeSlice[0]) {
		return invalidStringReturn()
	}

	length := len(runeSlice)

	if length == 1 {
		isBSRune := utils.IsBackSlashRune(runeSlice[0])
		isSpaceRune := utils.IsSpaceRune(runeSlice[0])
		wrongRune := isBSRune || isSpaceRune
		if wrongRune {
			return invalidStringReturn()
		}
		return passedStringReturn(s)
	}

	kGrow := 4
	context := newUnpackerContext(kGrow * length)

	for i := 0; i < length; i++ {
		r := runeSlice[i]

		// space characters are not allowed
		if utils.IsSpaceRune(r) {
			return invalidStringReturn()
		}

		isLastIter := i == length-1

		isDigitRune := utils.IsDigitRune(r)
		isBSRune := utils.IsBackSlashRune(r)

		isNumberMode := context.isNumberMode()
		isEscapeMode := context.isEscapeMode()

		if isEscapeMode {
			isValidEscapeRune := isDigitRune || isBSRune
			// only digits or backslashes are allowed for escape
			if !isValidEscapeRune {
				return invalidStringReturn()
			}
			context.writeRuneToResult(r)
			context.setEscapeModeOff()
			continue
		}

		if isBSRune && !isEscapeMode {
			// string with escape that escapes nothing considered as invalid
			if isLastIter {
				return invalidStringReturn()
			}

			if isNumberMode {
				if err := finishNumberUnpacking(context); err != nil {
					fmt.Fprintln(os.Stderr, err)
					return invalidStringReturn()
				}
			}
			context.setEscapeModeOn()
			continue
		}

		if isDigitRune && isNumberMode {
			context.writeRuneToNumber(r)
			continue
		}
		if isDigitRune && !isNumberMode {
			context.setUnpackingNumberOn()
			context.writeRuneToNumber(r)
			continue
		}
		if !isDigitRune && isNumberMode {
			if err := finishNumberUnpacking(context); err != nil {
				fmt.Fprintln(os.Stderr, err)
				return invalidStringReturn()
			}
			context.writeRuneToResult(r)
			continue
		}
		if !isDigitRune && !isNumberMode {
			context.writeRuneToResult(r)
			continue
		}
	}

	if context.isNumberMode() {
		if err := finishNumberUnpacking(context); err != nil {
			fmt.Fprintln(os.Stderr, err)
			return invalidStringReturn()
		}
	}

	return context.getResultString()
}

func finishNumberUnpacking(context *unpackerContext) error {
	n, err := context.getAccumulatedNumber()
	if err != nil {
		return err
	}
	lr := context.lastResultRune
	// we've wrote one rune to the result
	nWrite := n - 1
	context.writeDuplicatedRuneToResult(lr, nWrite)
	context.finishNumberUnpacking()
	return nil
}

func passedStringReturn(s string) (string, error) {
	return s, nil
}

func invalidStringReturn() (string, error) {
	return "", ErrInvalidString
}
