package resp

import (
	"bufio"
	"fmt"
)

func WriteSimple(writer *bufio.Writer, msg string) {
	writer.WriteString("+" + msg + "\r\n")
}

func WriteError(writer *bufio.Writer, msg string) {
	writer.WriteString("-ERR " + msg + "\r\n")
}

func WriteBulkString(writer *bufio.Writer, msg string) {
	writer.WriteString(fmt.Sprintf("$%d\r\n%s\r\n", len(msg), msg))
}

func WriteArrayBulk(writer *bufio.Writer, msgs ...string) {
	writer.WriteString(fmt.Sprintf("*%d\r\n", len(msgs)))
	for _, s := range msgs {
		writer.WriteString(fmt.Sprintf("$%d\r\n%s\r\n", len(s), s))
	}
}

func WriteNil(writer *bufio.Writer) {
	writer.WriteString("$-1\r\n")
}
