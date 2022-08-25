package command

import "fmt"

type cHelp struct {
	arg string
}

func newCommandHelp(args string) (*cHelp, error) {
	obj := &cHelp{}
	obj.arg = args
	return obj, nil
}
func (c *cHelp) Execute() error {
	c.printHelp()
	return nil
}

func (c *cHelp) printHelp() {
	fmt.Println("List of commands available:")
	fmt.Printf("%s: navigate through file system\n", cmdCD)
	fmt.Printf("%s: show info about current processes\n", cmdPS)
	fmt.Printf("%s: print working directory\n", cmdPWD)
	fmt.Printf("%s: print arguments passed\n", cmdECHO)
	fmt.Printf("%s: kill process by PID or name\n", cmdKILL)
	fmt.Printf("%s: fork current process\n", cmdFORK)
	fmt.Printf("%s: execute binary (might be used as %s& to run in forked process)\n", cmdEXEC, cmdEXEC)
	fmt.Printf("%s: quit from the shell\n", cmdQUIT)
}
