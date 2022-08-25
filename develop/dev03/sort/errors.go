package sort

import (
	"errors"
	"fmt"
)

var (
	errInvalidSortCol         = fmt.Errorf("sort: sort column number starts with %d", minSortCol)
	errBadOpenFile            = errors.New("sort: can't read file")
	errBadCreateFile          = errors.New("sort: can't create file")
	errWriteFail              = errors.New("sort: can't write to file")
	errStrNotStartsWithNumber = errors.New("sort: string not starts with number")
)
