package telnet

import (
	"errors"
	"flag"
	"net"
	"time"
)

const (
	defautlTimeout = 10 * time.Second
)

type options struct {
	address string
	timeout time.Duration
}

func newOptions(args []string) (*options, error) {
	opt := &options{}
	fs := flag.NewFlagSet("go-telnet", flag.ContinueOnError)
	fs.DurationVar(&opt.timeout, "timeout", defautlTimeout, "connection timeoout")

	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	host, port := fs.Arg(0), fs.Arg(1)
	if host == "" || port == "" {
		return nil, errors.New("invalid connetction address was provided")
	}

	opt.address = net.JoinHostPort(host, port)

	return opt, nil
}
