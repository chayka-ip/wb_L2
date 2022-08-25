package unpackstr

import (
	"fmt"
	"os"
)

//ExecuteCLI unpacks string passed as argument
func ExecuteCLI(args []string) int {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "string is not provided")
		return 2
	}

	str, err := unpack(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Runtime error: %v\n", err)
		return 1
	}

	fmt.Printf("Result: %s\n", str)

	return 0
}
