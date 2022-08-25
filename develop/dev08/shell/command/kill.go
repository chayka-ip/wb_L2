package command

import (
	"os"
	"strconv"
	"strings"

	"github.com/mitchellh/go-ps"
)

type cKill struct {
	nameOrPid string
}

func newCommandKill(args string) (*cKill, error) {
	pidOrName, err := parseArg(args)
	if err != nil {
		return nil, err
	}
	obj := &cKill{}
	obj.nameOrPid = pidOrName
	return obj, nil
}
func (c *cKill) Execute() error {
	target := c.nameOrPid
	pid, err := strconv.Atoi(c.nameOrPid)
	if err != nil {
		pid, err = getPIDByExecutableName(target)
		if err != nil {
			return err
		}
	}
	return killProcByPid(pid)
}

func killProcByPid(pid int) error {
	p, err := os.FindProcess(pid)
	if err != nil {
		return err
	}
	return p.Kill()
}

func parseArg(args string) (string, error) {
	a := strings.Fields(args)
	if len(a) == 1 {
		return a[0], nil
	}
	return "", ErrBadArgs
}

func getPIDByExecutableName(name string) (int, error) {
	procList, err := ps.Processes()
	if err != nil {
		return 0, err
	}
	for _, v := range procList {
		if v.Executable() == name {
			return v.Pid(), nil
		}
	}
	return 0, UnknownProcess(name)
}
