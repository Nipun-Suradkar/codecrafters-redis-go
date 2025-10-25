package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
)

// Ensures gofmt doesn't remove the "net" and "os" imports in stage 1 (feel free to remove this!)
var _ = net.Listen
var _ = os.Exit

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		go handleConnection(conn)
	}

}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)
	commandHandler := NewCommandHelper(writer)
	for {
		// Try to decode a RESP array from the reader
		cmd, err := decodeRESPFromReader(reader)
		if err != nil {
			if err == io.EOF {
				fmt.Println("Connection closed")
			} else {
				fmt.Println("Error decoding RESP:", err)
			}
			return
		}
		fmt.Println("Decoded command:", cmd)

		commandHandler.HandleCommand(cmd)

	}
}
