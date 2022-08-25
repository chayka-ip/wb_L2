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
	htmlExt = ".html"
)

//ExecuteCLI ...
func ExecuteCLI(args []string) int {
	opt, err := newOptions(args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 2
	}

	if parser.HasExtension(opt.url.AbsoluteURL) {
		if _, err := downloadResource(opt.url); err != nil {
			fmt.Println(err)
			return 1
		}
		return 0
	}

	if err := downloadPage(opt, opt.url, 0); err != nil {
		fmt.Println(err)
		return 1
	}
	return 0
}

func downloadPage(opt *options, url *ResourceURL, currentDepth int) error {
	if !strings.HasSuffix(url.ResourceFullPath, htmlExt) {
		url.ResourceFullPath += htmlExt
	}

	body, err := downloadResource(url)
	if err != nil {
		return err
	}

	s := string(body)
	linkContainer := parser.ExtractLinksFromString(s)

	// download resources
	if opt.prerequesits {
		resourceLinks := linkContainer.GetResourceLinks()
		for _, v := range resourceLinks {
			rURL, err := makeURLFromTypedString(v.Link, url.Schema, url.Host, opt.downloadDir, v.URLType)
			if err != nil {
				fmt.Println(err)
				continue
			}
			_, err = downloadResource(rURL)
			if err != nil {
				fmt.Println(err)
			}
		}
	}

	// download pages
	if currentDepth < opt.depth {
		pageLinks := linkContainer.GetPageLinks()
		for _, v := range pageLinks {
			rURL, err := makeURLFromTypedString(v.Link, url.Schema, url.Host, opt.downloadDir, v.URLType)
			if err != nil {
				fmt.Println(err)
				continue
			}
			err = downloadPage(opt, rURL, currentDepth+1)
			if err != nil {
				fmt.Println(err)
			}
		}
		parser.AddSuffixToPageLinks(linkContainer, htmlExt)
	}

	err = parser.ConvertLinksInFileToRelative(*linkContainer, url.ResourceFullPath, url.Host)
	if err != nil {
		fmt.Println(err)
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
