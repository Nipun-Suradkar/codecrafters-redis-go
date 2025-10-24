package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
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

		handleCommand(cmd, writer)

	}
}

func handleCommand(cmd []string, writer *bufio.Writer) {
	if len(cmd) == 0 {
		return
	}

	if strings.EqualFold(cmd[0], "echo") {
		if len(cmd) > 1 {
			resp := fmt.Sprintf("$%d\r\n%s\r\n", len(cmd[1]), cmd[1])
			writer.WriteString(resp)
			writer.Flush()
		} else {
			writer.WriteString("-ERR missing argument\r\n")
			writer.Flush()
		}
	} else if strings.EqualFold(cmd[0], "ping") {
		writer.WriteString("+PONG\r\n")
		writer.Flush()
	} else {
		writer.WriteString("-ERR unknown command\r\n")
		writer.Flush()
	}
}
