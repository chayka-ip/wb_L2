package command

import (
	"fmt"
	"os"
)

type cQuit struct {
}

func newCommandQuit(args string) (*cQuit, error) {
	if err := validateNumArgs(0, 0, args); err != nil {
		return nil, err
	}

	obj := &cQuit{}
	return obj, nil
}
func (c *cQuit) Execute() error {
	fmt.Println("Go Shell is terminating...")
	os.Exit(0)
	return nil
}
