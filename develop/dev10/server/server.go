package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
)

const (
	port = "1234"
)

func main() {
	addr := ":" + port
	listen, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	defer listen.Close()

	fmt.Println("Server is running on ", addr)

	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println(err)
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	addr := conn.RemoteAddr()
	fmt.Println("New connection: ", addr)
	defer conn.Close()

	for {
		s, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			if err == io.EOF {
				fmt.Println("Closing connection with ", addr)
				return
			}

			fmt.Println("unexpected error: ", err)
			return
		}

		out := fmt.Sprintf("server echo: %s\n", s)
		conn.Write([]byte(out))
	}
}
