package command

import (
	"os"
	"os/exec"
	"strings"
	"syscall"
)

type cExec struct {
	path string
	args []string
}

func newCommandExec(args string) (*cExec, error) {
	path, argSlice, err := getPathAndArgs(args)
	if err != nil {
		return nil, err
	}

	obj := &cExec{}
	obj.path = path
	obj.args = argSlice
	return obj, nil
}
func (c *cExec) Execute() error {
	path, err := exec.LookPath(c.path)
	if err != nil {
		return err
	}

	return syscall.Exec(path, c.args, os.Environ())
}

func getPathAndArgs(arg string) (string, []string, error) {
	s := strings.Fields(arg)
	if len(s) == 0 {
		return "", nil, ErrNoArgs
	}
	return s[0], s[1:], nil
}
