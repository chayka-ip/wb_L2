package command

import (
	"go-shell/shell/utils"
	"strings"
)

const (
	cmdCD   = "cd"
	cmdPS   = "ps"
	cmdPWD  = "pwd"
	cmdECHO = "echo"
	cmdKILL = "kill"
	cmdFORK = "fork"
	cmdEXEC = "exec"
	cmdHELP = "help"
	cmdQUIT = "q"

	forkAmp = "&"
)

//ICommand ...
type ICommand interface {
	Execute() error
}

//NewCommand creates new command
func NewCommand(rawArgs string) (ICommand, error) {
	command, args, err := extractCommandAndArgs(rawArgs)
	if err != nil {
		return nil, ErrBadArgs
	}
	switch command {
	case cmdCD:
		return newCommandCD(args)
	case cmdPS:
		return newCommandPs(args)
	case cmdPWD:
		return newCommandPwd(args)
	case cmdECHO:
		return newCommandEcho(args)
	case cmdKILL:
		return newCommandKill(args)
	case cmdEXEC:
		return newCommandExec(args)
	case cmdFORK:
		return newCommandFork(args)
	case cmdHELP:
		return newCommandHelp(args)
	case cmdQUIT:
		return newCommandQuit(args)
	default:
		return nil, UnknownCommand(command)
	}
}

func validateNumArgs(min, max int, args string) error {
	return utils.ValidateNumArgs(min, max, args, ErrTooFewArgs, ErrTooManyArgs)
}

func extractCommandAndArgs(args string) (string, string, error) {
	cmd, restArgs, err := utils.GetFirstArgAndRestAsStr(args)
	if err != nil {
		return "", "", err
	}

	isFork := strings.HasSuffix(cmd, forkAmp)
	if isFork {
		cmd = strings.TrimSuffix(cmd, forkAmp)
		restArgs = cmd + restArgs
		return cmdFORK, restArgs, nil
	}

	return cmd, restArgs, nil
}
