package command

import (
	"errors"
	"fmt"
)

var (
	//ErrNoArgs is returned when no arguments are provided
	ErrNoArgs = errors.New("no arguments provided")
	//ErrBadArgs is returned when provided arguments are incorret
	ErrBadArgs = errors.New("bad arguments provided")
	//ErrTooManyArgs is returned when there are too many arguments provided
	ErrTooManyArgs = errors.New("too many arguments")
	//ErrTooFewArgs is returned when there are too many arguments provided
	ErrTooFewArgs = errors.New("too few arguments")
)

//UnknownCommand is returned when requested command is unknown
func UnknownCommand(name string) error {
	return fmt.Errorf("invalid command: %s", name)
}

//UnknownProcess is returned when requested process is unknown
func UnknownProcess(name string) error {
	return fmt.Errorf("go-shell: process is unknown: %s", name)
}
