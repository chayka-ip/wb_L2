package wget

import (
	"fmt"
	"go-wget/wget/parser"
	"os"
	"strings"
)

func makeDirTree(path string) error {
	return os.MkdirAll(path, 0777)
}
func joinPathURL(s ...string) string {
	return strings.Join(s, slash)
}

func makeURLFromTypedString(rawURL, schema, host, downloadDir string, urlType parser.URLType) (*ResourceURL, error) {
	u := ""
	switch urlType {
	case parser.URLTypeAbsolute:
		u = rawURL
	case parser.URLTypeRelative:
		u = fmt.Sprintf("%s://%s/%s", schema, host, rawURL)
	case parser.URLTypeThisOtherDomainNoSchema:
		u = fmt.Sprintf("%s:%s", schema, rawURL)
	}
	return NewResourceURL(u, downloadDir)
}

func logDownload(url string, size int) {
	fmt.Printf("Got %d bytes from %s\n", size, url)
}
