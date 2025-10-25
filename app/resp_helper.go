package main

import (
	"bufio"
	"fmt"
	"strconv"
)

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
