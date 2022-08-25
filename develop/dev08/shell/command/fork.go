package command

import (
	"syscall"
)

// partially works but cant kill and

type cFork struct {
	args string
}

func newCommandFork(args string) (*cFork, error) {
	obj := &cFork{}
	obj.args = args
	return obj, nil
}
func (c *cFork) Execute() error {
	childPID, _, err := syscall.Syscall(syscall.SYS_FORK, 0, 0, 0)

	if err != 0 {
		return err
	}

	if childPID == 0 {
		cm, err := NewCommand(c.args)
		if err != nil {
			return err
		}
		return cm.Execute()
	}

	return nil
}
