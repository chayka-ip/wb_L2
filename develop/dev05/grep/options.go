package grep

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
)

var (
	errBadOpenFile = errors.New("grep: can't read file")
)

type optionsRaw struct {
	numLinesBeforeMatch int
	numLinesAfterMatch  int
	numContextLines     int

	isPrintMatchCount bool
	isIgnoreCase      bool
	isInvertSearch    bool
	isExactMatch      bool
	isPrintLineNum    bool
}

type options struct {
	pattern             string
	numLinesBeforeMatch int
	numLinesAfterMatch  int

	isPrintMatchCount bool
	isIgnoreCase      bool
	isInvertSearch    bool
	isExactMatch      bool
	isPrintLineNum    bool
	reader            io.ReadCloser
}

func newOptions(args []string) (*options, error) {
	optRaw := &optionsRaw{}
	fs := flag.NewFlagSet("flags", flag.ContinueOnError)
	fs.IntVar(&optRaw.numLinesAfterMatch, "A", 0, "show N lines after matched line")
	fs.IntVar(&optRaw.numLinesBeforeMatch, "B", 0, "show N lines before matched line")
	fs.IntVar(&optRaw.numContextLines, "C", 0, "show N lines before and N lines after matched line")
	fs.BoolVar(&optRaw.isPrintMatchCount, "c", false, "print matches count")
	fs.BoolVar(&optRaw.isInvertSearch, "v", false, "exclude lines that contain pattern")
	fs.BoolVar(&optRaw.isIgnoreCase, "i", false, "ignore case")
	fs.BoolVar(&optRaw.isExactMatch, "F", false, "check for exact matched entries instead of regexp pattern")
	fs.BoolVar(&optRaw.isPrintMatchCount, "n", false, "print line number before each line")

	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	rawPattern := fs.Arg(0)
	filepath := fs.Arg(1)

	pattern, err := getPattern(rawPattern, optRaw.isExactMatch, optRaw.isIgnoreCase)
	if err != nil {
		return nil, err
	}

	linesBefore, linesAfter := optRaw.getNumLinesBeforeAndAfter()
	reader, err := getReader(filepath)
	if err != nil {
		return nil, err
	}

	opt := &options{
		pattern:             pattern,
		numLinesBeforeMatch: linesBefore,
		numLinesAfterMatch:  linesAfter,

		isPrintMatchCount: optRaw.isPrintMatchCount,
		isIgnoreCase:      optRaw.isIgnoreCase,
		isInvertSearch:    optRaw.isInvertSearch,
		isExactMatch:      optRaw.isExactMatch,
		isPrintLineNum:    optRaw.isPrintLineNum,
		reader:            reader,
	}
	return opt, nil
}

func (o *options) GetExactPattern() string {
	return getExactPattern(o.pattern, o.isIgnoreCase)
}

func (o *options) GetRegExpPattern() *regexp.Regexp {
	r, _ := getRegexpPattern(o.pattern, o.isIgnoreCase)
	return regexp.MustCompile(r)
}

func (o *optionsRaw) getNumLinesBeforeAndAfter() (int, int) {
	before := maxInt(o.numLinesBeforeMatch, o.numContextLines)
	after := maxInt(o.numLinesAfterMatch, o.numContextLines)
	return before, after
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
