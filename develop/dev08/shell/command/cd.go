package command

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode"
)

const (
	parentDirCode   = ".."
	dirRelativePref = "./"
	dirHome         = "/"
	slash           = "/"
)

type cCd struct {
	dir string
}

func newCommandCD(args string) (*cCd, error) {
	dirCurrent, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	dir, err := getDirFromArgs(dirCurrent, args)
	if err != nil {
		return nil, err
	}

	obj := &cCd{}
	obj.dir = dir
	return obj, nil
}

func (c *cCd) Execute() error {
	dir := c.dir
	if err := os.Chdir(dir); err != nil {
		return err
	}
	fmt.Println(dir)
	return nil
}

func getDirFromArgs(currentDir, args string) (string, error) {
	dir, err := parseArgs(args)
	if err != nil {
		return "", err
	}
	return getDirStringFromArgs(currentDir, dir), nil
}

func getDirStringFromArgs(currentDir, dir string) string {
	if dir == parentDirCode {
		return filepath.Dir(currentDir)
	}

	if dir == dirHome {
		return dirHome
	}

	// relative dir case
	if strings.HasPrefix(dir, dirRelativePref) {
		dir = strings.Trim(dir, dirRelativePref)
		out := jouinDirsWithCorrectSlash(currentDir, dir)
		return out
	}

	// handle absolute case
	if strings.HasPrefix(dir, slash) {
		return dir
	}

	dir = strings.TrimLeft(dir, ".")
	if dir == "" {
		return currentDir
	}

	out := jouinDirsWithCorrectSlash(currentDir, dir)
	return out
}

func jouinDirsWithCorrectSlash(currentDir string, dir string) string {
	if !strings.HasSuffix(currentDir, slash) {
		return fmt.Sprintf("%s/%s", currentDir, dir)
	}
	return currentDir + dir
}

func parseArgs(s string) (string, error) {
	s = strings.TrimSpace(s)

	hasArg := false
	var arg strings.Builder
	arg.Grow(20)

	if len(s) == 0 {
		return "", nil
	}

	for _, r := range s {
		notSpace := !unicode.IsSpace(r)
		if notSpace {
			arg.WriteRune(r)
			hasArg = true
			continue
		}

		if notSpace && hasArg {
			return "", ErrTooManyArgs
		}
	}

	return arg.String(), nil
}
