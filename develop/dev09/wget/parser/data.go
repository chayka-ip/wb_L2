package parser

import (
	"strings"
)

//URLType ...
type URLType uint8

const (
	URLTypeAbsolute                URLType = 0
	URLTypeRelative                URLType = 1
	URLTypeThisOtherDomainNoSchema URLType = 2
	URLTypeUnknown                 URLType = 255
)

//LinkItem ...
type LinkItem struct {
	// byte start position in source chunk
	ByteStartPos    int
	OrigialDataSize int
	Link            string
	URLType         URLType
}

func (i *LinkItem) hasProperURLType() bool {
	return i.URLType != URLTypeUnknown
}

func (i *LinkItem) getNumBytesDelta() int {
	return len(i.Link) - i.OrigialDataSize
}

//LinkContainer ...
type LinkContainer struct {
	// all links found
	Data             []LinkItem
	PageLinksInd     []int
	ResourceLinksInd []int
}

//NewLinkContainer ...
func NewLinkContainer() *LinkContainer {
	return &LinkContainer{
		Data: make([]LinkItem, 0),
	}
}

func (p *LinkContainer) add(d LinkItem) {
	p.Data = append(p.Data, d)
}

//GetDataAsStringSlice ...
func (p *LinkContainer) GetDataAsStringSlice() []string {
	out := []string{}
	for _, v := range p.Data {
		out = append(out, v.Link)
	}
	return out
}

//GetLinksAsStringSlice ...
func (p *LinkContainer) GetLinksAsStringSlice(indices []int) []string {
	out := make([]string, 0, len(indices))
	for _, v := range indices {
		out = append(out, p.Data[v].Link)
	}
	return out
}

//GetPageLinksAsStringSlice ...
func (p *LinkContainer) GetPageLinksAsStringSlice() []string {
	return p.GetLinksAsStringSlice(p.PageLinksInd)
}

//GetResourceLinksAsStringSlice ...
func (p *LinkContainer) GetResourceLinksAsStringSlice() []string {
	return p.GetLinksAsStringSlice(p.ResourceLinksInd)
}

//GetLinks ...
func (p *LinkContainer) GetLinks(indices []int) []LinkItem {
	out := make([]LinkItem, 0, len(indices))
	for _, v := range indices {
		out = append(out, p.Data[v])
	}
	return out
}

//GetPageLinks ...
func (p *LinkContainer) GetPageLinks() []LinkItem {
	return p.GetLinks(p.PageLinksInd)
}

//GetResourceLinks ...
func (p *LinkContainer) GetResourceLinks() []LinkItem {
	return p.GetLinks(p.ResourceLinksInd)
}

func (p *LinkContainer) removeImproperLinks() {
	out := make([]LinkItem, 0, len(p.Data))
	for _, v := range p.Data {
		c1 := strings.Contains(v.Link, slash)
		c2 := v.hasProperURLType()

		ok := c1 && c2

		if ok {
			out = append(out, v)
		}
	}
	p.Data = out
}

func (p *LinkContainer) searchPageAndResourceLinks() {
	for i, v := range p.Data {
		if HasExtension(v.Link) {
			p.ResourceLinksInd = append(p.ResourceLinksInd, i)
			continue
		}
		p.PageLinksInd = append(p.PageLinksInd, i)
	}
}
