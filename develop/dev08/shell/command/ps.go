package command

import (
	"fmt"
	"strings"

	"github.com/mitchellh/go-ps"
)

type cPs struct {
}

func newCommandPs(args string) (*cPs, error) {
	if err := validateNumArgs(0, 0, args); err != nil {
		return nil, err
	}

	obj := &cPs{}
	return obj, nil
}
func (c *cPs) Execute() error {
	r, err := ps.Processes()
	if err != nil {
		return err
	}
	head := fmt.Sprintf("%7s | %s", "PID", "executable")

	fmt.Println(head)
	fmt.Println(strings.Repeat("-", len(head)+5))

	for _, v := range r {
		fmt.Printf("%7d | %s\n", v.Pid(), v.Executable())
	}
	return nil
}
