package command

import (
	"fmt"
	"strings"
)

type cEcho struct {
	arg string
}

func newCommandEcho(args string) (*cEcho, error) {
	args = strings.TrimLeft(args, " ")

	obj := &cEcho{}
	obj.arg = args
	return obj, nil
}
func (c *cEcho) Execute() error {
	fmt.Println(c.arg)
	return nil
}
