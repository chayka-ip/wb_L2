package shell

import (
	"bufio"
	"fmt"
	"go-shell/shell/command"
	"go-shell/shell/utils"
	"io"
	"os"
	"strings"
)

//ExecuteCLI runs shell command
func ExecuteCLI() int {
	sh, err := newShell()

	if err != nil {
		fmt.Println(err)
		return 2
	}

	if err := sh.run(); err != nil {
		fmt.Println(err)
		return 1
	}

	return 0
}

type shell struct {
	reader io.ReadCloser
}

func newShell() (*shell, error) {
	sh := &shell{
		reader: os.Stdin,
	}
	return sh, nil
}

func (s *shell) run() error {
	fmt.Println("Go Shell is ready...")

	// start polling

	scanner := bufio.NewScanner(s.reader)
	for scanner.Scan() {
		args := scanner.Text()

		if utils.HasPipe(args) {
			argsSlice := strings.Split(args, utils.PIPE)
			for _, v := range argsSlice {
				s.execute(v)
			}

		} else {
			s.execute(args)
		}

		if err := scanner.Err(); err != nil {
			logError(err)
			return err
		}
	}

	return nil
}

func (s *shell) execute(args string) {
	command, err := command.NewCommand(args)
	if err != nil {
		logError(err)
		return
	}

	if err := command.Execute(); err != nil {
		logError(err)
	}
}

func logError(e error) {
	fmt.Println(e)
}
