package cut

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"
)

var (
	//ErrBadFields is returned when ivalid field arg was passed
	ErrBadFields = fmt.Errorf("cut: fileds are numbered from %d to %d", minSearchRange, maxSearchRange)
	//ErrBadDelimeter is returned when ivalid delimeter arg was passed
	ErrBadDelimeter = errors.New("cut: the delimeter must be a single character")
	//ErrBadReader is returned when data reader can't be obtained
	ErrBadReader = errors.New("cut: no os reader is provided")
)

const (
	defaultDelimiter = string('\t')
	minSearchRange   = 1
	maxSearchRange   = 1000
)

type rawOptions struct {
	fields    string
	delimiter string
	separated bool
}

type options struct {
	fields      *searchRange
	delimiter   string
	withSepOnly bool
	reader      io.Reader
}

func newOptions(args []string) (*options, error) {
	optRaw := rawOptions{}

	fs := flag.NewFlagSet("options", flag.ExitOnError)
	fs.StringVar(&optRaw.fields, "f", "", "fields")
	fs.StringVar(&optRaw.delimiter, "d", defaultDelimiter, "delimiter")
	fs.BoolVar(&optRaw.separated, "s", false, "show lines with separator")

	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	delimeter, err := optRaw.getDelimiter()
	if err != nil {
		return nil, err
	}
	fields, err := parseFields(optRaw.fields)
	if err != nil {
		return nil, err
	}

	reader, err := getReader()
	if err != nil {
		return nil, err
	}

	opt := &options{
		fields:      fields,
		delimiter:   delimeter,
		withSepOnly: optRaw.separated,
		reader:      reader,
	}

	return opt, nil
}

func (o *rawOptions) getDelimiter() (string, error) {
	d := o.delimiter
	if utf8.RuneCountInString(d) == 1 {
		return d, nil
	}
	return "", ErrBadDelimeter
}

func getReader() (io.Reader, error) {
	if stat, err := os.Stdin.Stat(); err == nil {
		if (stat.Mode() & os.ModeCharDevice) == 0 {
			return os.Stdin, nil
		}
	}
	return nil, ErrBadReader
}

func parseFields(s string) (*searchRange, error) {
	result := newSearchRange()
	space := regexp.MustCompile(`\s+`)
	s = space.ReplaceAllString(s, "")

	if s == "" {
		result.addRange(minSearchRange, maxSearchRange)
		return result, nil
	}

	p := strings.Split(s, ",")

	for _, chunk := range p {
		num, err := parseFieldValue(chunk)
		ok := (err == nil) && (num >= minSearchRange)
		if ok {
			result.addSingle(num)
			continue
		}
		rMin, rMax, err := unpackIntRangeEnds(chunk)
		if err != nil {
			return nil, err
		}
		result.addRange(rMin, rMax)
	}

	return result, nil
}

//Returns min and max bounds of the interval (inclusive)
//Returns error if a < b.
//Expects range in form "a-b".
//If a is ommited - a = minSearchRange.
//If b is ommited - b = maxSearchRange.
func unpackIntRangeEnds(s string) (int, int, error) {
	badarg := func() (int, int, error) { return 0, 0, ErrBadFields }
	parseV := func(v string, defaultVal int) (int, error) {
		if v == "" {
			return defaultVal, nil
		}
		return parseFieldValue(v)
	}

	p := strings.Split(s, "-")
	if len(p) != 2 {
		return badarg()
	}

	if (p[0] == p[1]) && (p[0] == "") {
		return badarg()
	}

	a, err1 := parseV(p[0], minSearchRange)
	b, err2 := parseV(p[1], maxSearchRange)

	if (err1 != nil) || (err2 != nil) || a > b {
		return badarg()
	}

	return a, b, nil
}

func parseFieldValue(v string) (int, error) {
	d, err := strconv.Atoi(v)
	if err != nil {
		return 0, err
	}

	bad := d < minSearchRange || d > maxSearchRange
	if bad {
		return 0, ErrBadFields
	}

	return d, nil
}
