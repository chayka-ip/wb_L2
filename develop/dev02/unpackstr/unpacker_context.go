package unpackstr

import (
	"bytes"
	"errors"
	"strconv"
	"unpack-str/utils"
)

var (
	errEmptyBuffer = errors.New("buffer is empty")
)

const (
	unpackModeNone   uint8 = 0
	unpackModeNumber uint8 = 1
	unpackModeEscape uint8 = 2
)

type unpackerContext struct {
	// result buffer
	buf            *bytes.Buffer
	lastResultRune rune
	// buffer to store number during unpacking
	numBuf *bytes.Buffer
	unpackMode uint8
}

func newUnpackerContext(growSize int) *unpackerContext {
	buf := &bytes.Buffer{}
	buf.Grow(growSize)
	numBuf := &bytes.Buffer{}
	return &unpackerContext{
		buf: buf,
		numBuf:     numBuf,
		unpackMode: unpackModeNone,
	}
}

func (c *unpackerContext) writeRuneToResult(r rune) {
	c.buf.WriteRune(r)
	c.lastResultRune = r
}

func (c *unpackerContext) writeDuplicatedRuneToResult(r rune, count int) {
	for i := 0; i < count; i++ {
		c.writeRuneToResult(r)
	}
}

//Writes rune representing digit to number slice.
//Panics if non-digit rune is passed.
func (c *unpackerContext) writeRuneToNumber(r rune) {
	if !utils.IsDigitRune(r) {
		panic("not number is written into number slice")
	}
	c.numBuf.WriteRune(r)
}

func (c *unpackerContext) setUnpackingNumberOn() {
	c.unpackMode = unpackModeNumber
}

func (c *unpackerContext) setUnpackingNumberOff() {
	c.unpackMode = unpackModeNone
}

func (c *unpackerContext) setEscapeModeOn() {
	c.unpackMode = unpackModeEscape
}

func (c *unpackerContext) setEscapeModeOff() {
	c.unpackMode = unpackModeNone
}

func (c *unpackerContext) isNumberMode() bool {
	return c.unpackMode == unpackModeNumber
}

func (c *unpackerContext) isEscapeMode() bool {
	return c.unpackMode == unpackModeEscape
}

func (c *unpackerContext) finishNumberUnpacking() {
	c.setUnpackingNumberOff()
	c.numBuf.Reset()
}

func (c *unpackerContext) getAccumulatedNumber() (int, error) {
	if c.numBuf.Len() == 0 {
		return 0, errEmptyBuffer
	}
	return strconv.Atoi(c.numBuf.String())
}

func (c *unpackerContext) getLastResultRune() (rune, error) {
	if c.buf.Len() == 0 {
		return 0, errEmptyBuffer
	}
	return c.lastResultRune, nil
}

func (c *unpackerContext) getResultString() (string, error) {
	if c.buf.Len() == 0 {
		return "", errEmptyBuffer
	}
	return c.buf.String(), nil
}
