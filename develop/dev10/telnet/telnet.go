package telnet

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

const (
	newLine = '\n'
)

//ExecuteCLI ...
func ExecuteCLI(args []string) int {
	setupSignals()

	opt, err := newOptions(args)
	if err != nil {
		fmt.Println(err)
		return 2
	}

	conn, err := connect(opt)
	if err != nil {
		fmt.Println(err)
		return 1
	}

	if err := handleConnection(opt, conn); err != nil {
		fmt.Println(err)
		return 1
	}

	return 0
}

func connect(opt *options) (net.Conn, error) {
	conn, err := net.DialTimeout("tcp", opt.address, opt.timeout)
	if err != nil {
		return nil, err
	}
	fmt.Println("Connection is established with ", opt.address)
	return conn, nil
}

func post(conn net.Conn, exitChan chan<- struct{}, errorChan chan<- error) {
	for {
		reader := bufio.NewReader(os.Stdin)
		s, err := reader.ReadString(newLine)
		if err != nil {
			handleErr(err, exitChan, errorChan)
			return
		}
		fmt.Fprintf(conn, s)
	}
}

func get(conn net.Conn, exitChan chan<- struct{}, errorChan chan<- error) {
	for {
		response, err := bufio.NewReader(conn).ReadString(newLine)
		if err != nil {
			handleErr(err, exitChan, errorChan)
			return
		}
		response = strings.TrimRight(response, string(newLine))
		fmt.Println("->", response)
	}
}

func handleErr(err error, exitChan chan<- struct{}, errorChan chan<- error) {
	if err != nil {
		// tracking for exit on Ctrl + D
		if err == io.EOF {
			exitChan <- struct{}{}
			return
		}
		errorChan <- err
		return
	}
}

func handleConnection(opt *options, conn net.Conn) error {
	defer conn.Close()
	exitChan := make(chan struct{})
	errorChan := make(chan error)

	go post(conn, exitChan, errorChan)
	go get(conn, exitChan, errorChan)

	select {
	case <-exitChan:
		fmt.Println("Terminating session...")
		time.Sleep(opt.timeout)
		fmt.Println("Connection is closed.")
		return nil
	case err := <-errorChan:
		return err
	}
}

func setupSignals() {
	signal.Ignore(syscall.SIGINT)
}
