package wget

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
)

const (
	slash         = "/"
	defautlScheme = "https"
)

//ResourceURL ...
type ResourceURL struct {
	// remote url
	AbsoluteURL string
	// remote url without schema
	URLNoSchema string
	// full path for file (if stored)
	ResourceFullPath string
	// path to dir for file (if stored)
	ResourceDirPath string
	Schema          string
	Host            string
}

//NewResourceURL ...
func NewResourceURL(absURL, downloadDir string) (*ResourceURL, error) {
	b, err := url.Parse(absURL)
	if err != nil {
		return nil, err
	}

	uAbs := b.String()
	sParts := strings.Split(absURL, "://")
	if len(sParts) < 2 {
		return nil, errors.New("url is invalid")
	}

	schema := sParts[0]
	urlNoSchema := strings.Split(absURL, "://")[1]
	rFullPath := fmt.Sprintf("%s/%s", downloadDir, urlNoSchema)
	parts := strings.Split(rFullPath, slash)
	host := parts[1]
	rDirPath := strings.Join(parts[:len(parts)-1], slash)

	u := &ResourceURL{
		AbsoluteURL:      uAbs,
		URLNoSchema:      urlNoSchema,
		ResourceFullPath: rFullPath,
		ResourceDirPath:  rDirPath,
		Schema:           schema,
		Host:             host,
	}
	return u, nil
}
