package parser

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

//HasExtension
func HasExtension(s string) bool {
	p := strings.Split(s, slash)
	lastEl := p[len(p)-1]
	if strings.Contains(lastEl, dot) {
		parts := strings.Split(lastEl, dot)
		ext := strings.ToLower(parts[len(parts)-1])
		return IsResourceFormat(ext)
	}
	return false
}

func linkToLocalResourceLinks(rawURL, host, rootDir string, urlType URLType) string {
	u := ""
	switch urlType {
	case URLTypeAbsolute:
		{
			p := strings.Split(rawURL, urlSchemaPrefix)
			u += strings.Join(p[1:], "")
		}
	case URLTypeRelative:
		rawURL = strings.TrimLeft(rawURL, slash)
		u += fmt.Sprintf("%s%s", host, rawURL)
	case URLTypeThisOtherDomainNoSchema:
		rawURL = strings.TrimLeft(rawURL, slash)
		u += rawURL
	}
	return fmt.Sprintf("%s/%s", rootDir, u)
}

//getNumLevelsToCommonParentDir calculates number of levels to up from pathRef
//to reach common directory
//Returns error if there are no common directory above,
//or files aren't in the same directory
func getNumLevelsToCommonParentDir(pathRef, pathTarget string) (int, error) {
	badReturn := func() (int, error) { return -1, errors.New("no common parent directory") }
	pr := strings.Split(pathRef, slash)
	pt := strings.Split(pathTarget, slash)

	nr, nt := len(pr), len(pt)
	if nr == 0 || nt == 0 {
		return badReturn()
	}

	if pathRef == pathTarget {
		return 0, nil
	}

	if pr[0] != pt[0] {
		return badReturn()
	}

	i := 0
	for i < nr {
		if i < nt {
			if pr[i] == pt[i] {
				i++
				continue
			}
		}
		break
	}
	return nr - i, nil
}

func resourceLinkToLocalRelative(targetPath, resourcePath string) (string, error) {
	relDir, err := getNumLevelsToCommonParentDir(targetPath, resourcePath)
	if err != nil {
		return "", err
	}
	if relDir == 0 {
		rs := slash + strings.TrimLeft(resourcePath, targetPath)
		return rs, nil
	}

	p := strings.Split(targetPath, slash)[:relDir]
	common := strings.Join(p, slash)
	rs := strings.TrimLeft(resourcePath, common)

	b := strings.Builder{}
	b.Grow(2 * relDir)
	for i := 0; i < relDir; i++ {
		b.WriteString(relativeDirPlaceHolder)
		b.WriteString(slash)
	}
	b.WriteString(rs)
	return b.String(), nil
}

func runeSliceEndsWith(s []rune, seq []rune) bool {
	sLen := len(s)
	seqLen := len(seq)
	if sLen < seqLen {
		return false
	}
	offset := sLen - seqLen
	for i := offset; i < sLen; i++ {
		if s[i] != seq[i-offset] {
			return false
		}
	}
	return true
}

func hasResource(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !fileInfo.IsDir()
}

//overrideBytes overrides bytes in src  starting from startPos
//if byteDelta > 0 this count will be added to src after startPos to fit new data
//if byteDelta < 0 redundant bytes will be deleted after startPos
func overrideBytes(src []byte, data []byte, startPos int, byteDelta int) ([]byte, error) {
	if startPos > len(src)-1 {
		return nil, errors.New("incorrect start position")
	}

	afterInd := startPos + len(data)
	endInd := afterInd - byteDelta
	s := make([]byte, len(src)+byteDelta)
	copy(s[:startPos], src[:startPos])
	copy(s[startPos:], data)
	copy(s[afterInd:], src[endInd:])

	return s, nil
}

//convertLinksToRelative tries to find corresponding file and modify links if it was found
//pageDirectory is path to directory where file to replace link is stored
func convertLinksToRelative(p *LinkContainer, pageDir, host, rootDir string) {
	byteOffset := 0

	md := make([]LinkItem, 0, len(p.Data))

	for _, v := range p.Data {
		v.ByteStartPos += byteOffset

		resourceLink := linkToLocalResourceLinks(v.Link, host, rootDir, v.URLType)

		if hasResource(resourceLink) {
			modifiedLink, err := resourceLinkToLocalRelative(pageDir, resourceLink)
			if err != nil {
				md = append(md, v)
				continue
			}
			byteOffset += len(modifiedLink) - len(v.Link)
			v.Link = modifiedLink
			md = append(md, v)
			continue
		}
		md = append(md, v)
	}
	p.Data = md
}

//overrideLinks replaces links in original content with contained data
//src is expected to be oriignal file on which this container is based
//does not modify src data
func overrideLinks(p *LinkContainer, src []byte) ([]byte, error) {
	out := append([]byte{}, src...)

	for _, v := range p.Data {
		n := v.getNumBytesDelta()
		data := []byte(v.Link)
		b, err := overrideBytes(out, data, v.ByteStartPos, n)
		if err != nil {
			return nil, err
		}
		out = b
	}
	return out, nil
}

//ConvertLinksInFileToRelative overrides provided file with relative links
func ConvertLinksInFileToRelative(p LinkContainer, filePath, host string) error {
	parts := strings.Split(filePath, slash)
	if len(parts) == 0 {
		return errors.New("invalid file path")
	}
	rootDir := parts[0]
	pathToFileDir := strings.Join(parts[:len(parts)-1], slash)

	convertLinksToRelative(&p, pathToFileDir, host, rootDir)
	fData, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	mData, err := overrideLinks(&p, fData)
	if err != nil {
		return err
	}

	err = os.WriteFile(filePath, mData, 0777)
	if err != nil {
		return err
	}

	return nil
}

//AddSuffixToPageLinks ...
func AddSuffixToPageLinks(p *LinkContainer, suffix string) {
	d := make([]LinkItem, 0, len(p.Data))

	pageIndices := make(map[int]struct{}, len(p.PageLinksInd))
	for _, ind := range p.PageLinksInd {
		pageIndices[ind] = struct{}{}
	}

	for i, v := range p.Data {
		if _, has := pageIndices[i]; has {
			if !strings.HasSuffix(v.Link, suffix) {
				v.Link += suffix
			}
		}
		d = append(d, v)
	}
	p.Data = d
}
