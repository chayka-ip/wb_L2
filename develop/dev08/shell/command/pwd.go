package command

import (
	"fmt"
	"os"
)

type cPwd struct {
}

func newCommandPwd(args string) (*cPwd, error) {
	if err := validateNumArgs(0, 0, args); err != nil {
		return nil, err
	}

	obj := &cPwd{}
	return obj, nil
}
func (c *cPwd) Execute() error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}
	fmt.Println(dir)
	return nil
}
