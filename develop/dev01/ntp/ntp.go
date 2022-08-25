package ntp

import (
	"flag"
	"fmt"
	"time"

	"github.com/beevik/ntp"
)

const (
	defaultHostNTP = "0.beevik-ntp.pool.ntp.org"
)

type options struct {
	//remote NTP server address
	host string
	//whether to format output time
	unixf bool
}

//ExecuteCLI is an entrypoint to NTP command from CLI.
//Returns exist status code after execution.
func ExecuteCLI(argsCLI []string) int {
	opt, err := parseArgs(argsCLI)
	if err != nil {
		fmt.Println(err)
		return 2
	}

	if err := printTime(opt); err != nil {
		fmt.Println(err)
		return 1
	}

	return 0
}

func parseArgs(argsCLI []string) (*options, error) {
	opt := &options{}
	fs := flag.NewFlagSet("ntp", flag.ContinueOnError)
	fs.StringVar(&opt.host, "host", defaultHostNTP, "ntp host")
	fs.BoolVar(&opt.unixf, "unixf", false, "apply unix date time format")
	if err := fs.Parse(argsCLI); err != nil {
		return nil, err
	}
	return opt, nil
}

func printTime(opt *options) error {
	ct, err := ntp.Time(opt.host)
	if err != nil {
		fmt.Println(err)
		return err
	}
	if opt.unixf {
		fmt.Println(ct.UTC().Format(time.UnixDate))
	} else {
		fmt.Println(ct)
	}
	return nil
}
