package wget

import (
	"fmt"
	"go-wget/wget/parser"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

const (
	pageSuffix = ".html"
)

type pages struct {
	urls []*ResourceURL
}

func newPages() *pages {
	return &pages{
		make([]*ResourceURL, 0),
	}
}

//ExecuteCLI ...
func ExecuteCLI(args []string) int {
	opt, err := newOptions(args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 2
	}

	if parser.HasResourceExtension(opt.url.AbsoluteURL) {
		if _, err := downloadResource(opt.url); err != nil {
			fmt.Println(err)
			return 1
		}
		return 0
	}

	pages := newPages()

	if err := downloadPage(opt, pages, opt.url, 0); err != nil {
		fmt.Println(err)
		return 1
	}

	convertLinksInPages(pages)
	fmt.Println("===== DONE =====")
	return 0
}

func downloadPage(opt *options, pages *pages, url *ResourceURL, currentDepth int) error {
	pages.urls = append(pages.urls, url)
	url.ResourceFullPath = strings.TrimRight(url.ResourceFullPath, slash)
	if !strings.HasSuffix(url.ResourceFullPath, pageSuffix) {
		url.ResourceFullPath += pageSuffix
	}

	body, err := downloadResource(url)
	if err != nil {
		return err
	}

	s := string(body)
	linkContainer := parser.ExtractLinksFromString(s)

	getURLS := func(p []parser.LinkItem) <-chan *ResourceURL {
		ch := make(chan *ResourceURL)
		go func() {
			for _, v := range p {
				rURL, err := makeURLFromTypedString(v.Link, url.Schema, url.Host, opt.downloadDir, v.URLType)
				if err != nil {
					fmt.Println(err)
					continue
				}
				ch <- rURL
			}
			close(ch)
		}()
		return ch
	}

	// download resources
	if opt.prerequesits {
		for v := range getURLS(linkContainer.GetResourceLinks()) {
			_, err = downloadResource(v)
			if err != nil {
				fmt.Println(err)
			}
		}
	}

	// download pages
	if currentDepth < opt.depth {
		for v := range getURLS(linkContainer.GetPageLinks()) {
			err = downloadPage(opt, pages, v, currentDepth+1)
			if err != nil {
				fmt.Println(err)
			}
		}
	}

	return nil
}

func downloadResource(url *ResourceURL) ([]byte, error) {
	resp, err := http.Get(url.AbsoluteURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	makeDirTree(url.ResourceDirPath)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	f, err := os.Create(url.ResourceFullPath)
	if err != nil {
		return nil, err
	}
	_, err = f.Write(body)
	if err != nil {
		return nil, err
	}

	logDownload(url.AbsoluteURL, len(body))
	return body, nil
}

func convertLinksInPages(pages *pages) {
	for _, url := range pages.urls {
		data, err := os.ReadFile(url.ResourceFullPath)
		if err != nil {
			fmt.Println(err)
			continue
		}
		linkContainer := parser.ExtractLinksFromString(string(data))
		err = parser.ConvertLinksInFileToRelative(*linkContainer, url.ResourceFullPath, url.Host, pageSuffix)
		if err != nil {
			fmt.Println(err)
		}
	}
}
