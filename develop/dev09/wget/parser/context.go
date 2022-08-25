package parser

type contextPars struct {
	//context accumulation
	ctx []rune
	// data accumulation
	data []rune
	// is extracting data
	isIncontext bool
	// where the actual data starts in source stream
	dataStartIndex int
	// sequences from wich context might start
	contextPrefixes [][]rune
	// sequence to end context
	contextEnd rune
}

func newContext(contextPrefixes []string, contextEnd rune) *contextPars {
	prefixes := [][]rune{}
	for _, v := range contextPrefixes {
		prefixes = append(prefixes, []rune(v))
	}

	return &contextPars{
		ctx:             make([]rune, 0),
		data:            make([]rune, 0),
		contextPrefixes: prefixes,
		contextEnd:      contextEnd,
	}
}

func (c *contextPars) startContext(startIndex int) {
	c.isIncontext = true
	c.dataStartIndex = startIndex
}

func (c *contextPars) reset() {
	c.ctx = make([]rune, 0)
	c.data = make([]rune, 0)
	c.isIncontext = false
	c.dataStartIndex = 0
}

func (c *contextPars) writeData(r rune) {
	c.data = append(c.data, r)
}

func (c *contextPars) writeContext(r rune) {
	c.ctx = append(c.ctx, r)
}
func (c *contextPars) write(r rune) {
	if c.isIncontext {
		c.writeData(r)
		return
	}
	c.writeContext(r)
}

func (c *contextPars) canStartContext() bool {
	for _, seq := range c.contextPrefixes {
		if runeSliceEndsWith(c.ctx, seq) {
			return true
		}
	}
	return false
}

func (c *contextPars) canEndContext() bool {
	if len(c.data) > 0 {
		return c.data[len(c.data)-1] == c.contextEnd
	}
	return false
}

func (c *contextPars) getDataAsString() string {
	ctxTrimInd := len(c.data) - 1
	return string(c.data[:ctxTrimInd])
}
