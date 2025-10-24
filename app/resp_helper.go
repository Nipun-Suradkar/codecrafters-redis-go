package main

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
)

var commands = map[string]func([]string, *bufio.Writer){
	"echo": func(cmd []string, writer *bufio.Writer) {
		if len(cmd) > 1 {
			resp := fmt.Sprintf("$%d\r\n%s\r\n", len(cmd[1]), cmd[1])
			writer.WriteString(resp)
			writer.Flush()
		} else {
			writer.WriteString("-ERR missing argument\r\n")
			writer.Flush()
		}
	},
	"ping": func(cmd []string, writer *bufio.Writer) {
		writer.WriteString("+PONG\r\n")
		writer.Flush()
	},
}

func decodeRESPFromReader(reader *bufio.Reader) ([]string, error) {
	// 1. Read first line (should be '*<numElements>\r\n')
	line, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}
	if len(line) < 2 || line[0] != '*' {
		return nil, fmt.Errorf("invalid RESP array header: %s", line)
	}

	numElements, err := strconv.Atoi(line[1 : len(line)-2]) // trim '*', '\r\n'
	if err != nil {
		return nil, fmt.Errorf("invalid number of elements: %v", err)
	}

	result := make([]string, 0, numElements)

	// 2. Read each bulk string: $<length>\r\n<data>\r\n
	for i := 0; i < numElements; i++ {
		// Read bulk string header
		header, err := reader.ReadString('\n')
		if err != nil {
			return nil, err
		}
		if len(header) < 2 || header[0] != '$' {
			return nil, fmt.Errorf("invalid bulk string header: %s", header)
		}

		strLen, err := strconv.Atoi(header[1 : len(header)-2])
		if err != nil {
			return nil, fmt.Errorf("invalid bulk string length: %v", err)
		}

		// Read bulk string data
		data := make([]byte, strLen+2) // include \r\n
		_, err = reader.Read(data)
		if err != nil {
			return nil, err
		}

		result = append(result, string(data[:strLen]))
	}
	return result, nil
}

func handleCommand(cmd []string, writer *bufio.Writer) {
	if len(cmd) == 0 {
		return
	}

	// Case-insensitive lookup
	for key, handler := range commands {
		if strings.EqualFold(cmd[0], key) {
			handler(cmd, writer)
			return
		}
	}

	// Default unknown command
	writer.WriteString("-ERR unknown command\r\n")
	writer.Flush()
}
