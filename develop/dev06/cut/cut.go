package cut

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

//ExecuteCLI executes cut command
func ExecuteCLI(args []string) int {
	opt, err := newOptions(args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 2
	}

	if err := run(opt); err != nil {
		fmt.Fprintf(os.Stderr, "Runtime error: %s\n", err)
		return 1
	}

	return 0
}

func run(opt *options) error {
	scanner := bufio.NewScanner(opt.reader)
	for scanner.Scan() {
		processLine(opt, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

func processLine(opt *options, s string) {
	chunks := strings.Split(s, opt.delimiter)
	num := len(chunks)

	if num == 1 {
		if !opt.withSepOnly {
			fmt.Println(chunks[0])
		}
		return
	}

	var parts []string
	for i := 0; i < num; i++ {
		if opt.fields.isInRange(i + 1) {
			parts = append(parts, chunks[i])
		}
	}
	if len(parts) > 0 {
		s := strings.Join(parts, opt.delimiter)
		fmt.Println(s)
	}
}
