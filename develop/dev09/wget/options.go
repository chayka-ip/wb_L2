package wget

import (
	"flag"
	"net/url"
)

type options struct {
	url          *ResourceURL
	depth        int
	prerequesits bool
	downloadDir  string
}

func newOptions(args []string) (*options, error) {
	opt := &options{}

	fs := flag.NewFlagSet("wget", flag.ContinueOnError)
	fs.IntVar(&opt.depth, "l", 0, "recursive download depth (from root=target_url)")
	fs.BoolVar(&opt.prerequesits, "p", false, "download addiional content (images/styles/etc)")
	fs.StringVar(&opt.downloadDir, "O", "data", "download directory name")

	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	url, err := url.Parse(fs.Arg(0))
	if err != nil {
		return nil, err
	}

	rURL, err := NewResourceURL(url.String(), opt.downloadDir)
	if err != nil {
		return nil, err
	}

	opt.url = rURL
	return opt, nil
}
