package parser

import (
	"strings"
	"unicode"
)

var (
	linkAttrs = []string{"href", "src", "content"}
)

const (
	urlSchemaPrefix              = "://"
	urlOtherDomainNoSchemaPrefix = "//"
	slash                        = "/"
	stringEndSuff                = '"'
	dot                          = "."
	relativeDirPlaceHolder       = ".."
)

//ExtractLinksFromString ...
func ExtractLinksFromString(s string) *LinkContainer {
	out := NewLinkContainer()
	prefixes := getLinkAttrsRaw()
	context := newContext(prefixes, stringEndSuff)

	for byteInd, r := range s {
		if context.isIncontext {
			if context.canEndContext() {
				item := finishContextAndGetData(context)
				out.add(item)
				continue
			}

			context.write(r)
			continue
		}

		if unicode.IsSpace(r) {
			continue
		}

		if context.canStartContext() {
			context.startContext(byteInd)
		}
		context.write(r)
	}

	if context.canEndContext() {
		item := finishContextAndGetData(context)
		out.add(item)
	}

	out.removeImproperLinks()
	out.searchPageAndResourceLinks()
	return out
}

func finishContextAndGetData(context *contextPars) LinkItem {
	link := context.getDataAsString()
	index := context.dataStartIndex
	urlType := getURLType(link)

	item := LinkItem{
		ByteStartPos:    index,
		Link:            link,
		URLType:         urlType,
		OrigialDataSize: len(link),
	}

	context.reset()
	return item
}

func getLinkAttrsRaw() []string {
	out := []string{}
	s := `="`
	for _, v := range linkAttrs {
		out = append(out, v+s)
	}
	return out
}

func getURLType(url string) URLType {
	if strings.HasPrefix(url, urlOtherDomainNoSchemaPrefix) {
		return URLTypeThisOtherDomainNoSchema
	}
	if strings.HasPrefix(url, slash) {
		return URLTypeRelative
	}
	if strings.Contains(url, urlSchemaPrefix) {
		return URLTypeAbsolute
	}
	return URLTypeUnknown
}
