package sort

import (
	"flag"
	"fmt"
	"io"
	"os"
)

const (
	defaultSeparator = " "

	minSortCol         = 1
	sortColDisabledVal = -1

	numericSortFlag      = "n"
	numericHumanSortFlag = "h"
	monthNameSortFlag    = "M"

	sortTypeDefautl      = 0
	sortTypeNumeric      = 1
	sortTypeHumanNumeric = 2
	sortTypeMonth        = 3
)

type flagOpt struct {
	value bool
	name  string
}

type optionsRaw struct {
	sortCol          int
	isNumericSort    bool
	isNumericSufSort bool
	isDescOrder      bool
	isMonthNameSort  bool
	isIgnoreTailBsp  bool
	isCheckIfSorted  bool
	isUniqueOnly     bool
}

type options struct {
	sortColIndex    int
	sortColEnabled  bool
	sortType        int
	isDescOrder     bool
	isCheckIfSorted bool
	isUniqueOnly    bool
	isIgnoreTailBsp bool
	separator       string
	reader          io.ReadCloser
	outFilePath     string
}

func newOptions(args []string) (*options, error) {
	optRaw := &optionsRaw{}
	fs := flag.NewFlagSet("flag", flag.ContinueOnError)
	fs.IntVar(&optRaw.sortCol, "k", minSortCol, "number of column to sort")
	fs.BoolVar(&optRaw.isNumericSort, numericSortFlag, false, "numeric sort")
	// fs.BoolVar(&optRaw.isNumericSufSort, numericHumanSortFlag, false, "numeric sort with suffix")
	// fs.BoolVar(&optRaw.isMonthNameSort, monthNameSortFlag, false, "sort by month name")
	fs.BoolVar(&optRaw.isDescOrder, "r", false, "sort in descending order")
	fs.BoolVar(&optRaw.isUniqueOnly, "u", false, "preserve only unique strings")
	// fs.BoolVar(&optRaw.isIgnoreTailBsp, "b", false, "ignore tail spaces")
	// fs.BoolVar(&optRaw.isCheckIfSorted, "c", false, "check if data is sorted")

	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	sortCol, sortColEnabled, err := optRaw.getSortColumn()
	if err != nil {
		return nil, err
	}

	if err := optRaw.validateIncompatibleOptions(); err != nil {
		return nil, err
	}

	reader, err := getReader(fs.Arg(0))
	if err != nil {
		return nil, err
	}

	outFilePath := fs.Arg(1)

	opt := &options{
		sortColIndex:    sortCol,
		sortColEnabled:  sortColEnabled,
		sortType:        optRaw.getSortType(),
		isDescOrder:     optRaw.isDescOrder,
		isCheckIfSorted: optRaw.isCheckIfSorted,
		isUniqueOnly:    optRaw.isUniqueOnly,
		isIgnoreTailBsp: optRaw.isIgnoreTailBsp,
		separator:       defaultSeparator,
		reader:          reader,
		outFilePath:     outFilePath,
	}

	return opt, nil
}

// returns real sort col (first means 0)
func (o *optionsRaw) getSortColumn() (int, bool, error) {
	if o.sortCol == sortColDisabledVal {
		return sortColDisabledVal, false, nil
	}
	if o.sortCol < minSortCol {
		return 0, false, errInvalidSortCol
	}
	return o.sortCol - 1, true, nil
}

func (o *optionsRaw) validateIncompatibleOptions() error {
	incompatible := []flagOpt{
		{o.isNumericSort, numericSortFlag},
		{o.isNumericSufSort, numericHumanSortFlag},
		{o.isMonthNameSort, monthNameSortFlag},
	}

	name1 := ""
	for _, f := range incompatible {
		if f.value {
			if name1 == "" {
				name1 = f.name
				continue
			}
			return o.incompatibleOptionsErr(name1, f.name)
		}
	}

	return nil
}

func (o *optionsRaw) getSortType() int {
	if o.isNumericSort {
		return sortTypeNumeric
	}
	if o.isNumericSufSort {
		return sortTypeHumanNumeric
	}
	if o.isMonthNameSort {
		return sortTypeMonth
	}
	return sortTypeDefautl
}

func getReader(filePath string) (io.ReadCloser, error) {
	if stat, err := os.Stdin.Stat(); err == nil {
		if (stat.Mode() & os.ModeCharDevice) == 0 {
			return os.Stdin, nil
		}
	}

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to open file %s: %v\n", filePath, err)
		return nil, errBadOpenFile
	}

	return file, nil
}

func (o *optionsRaw) incompatibleOptionsErr(a, b string) error {
	return fmt.Errorf("sort: options are incompatible: %s, %s", a, b)
}
