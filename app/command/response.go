package command

import "fmt"

func (c *Command) writeSimple(msg string) {
	c.Writer.WriteString("+" + msg + "\r\n")
}

func (c *Command) writeError(msg string) {
	c.Writer.WriteString("-ERR " + msg + "\r\n")
}

func (c *Command) writeBulkString(msg string) {
	c.Writer.WriteString(fmt.Sprintf("$%d\r\n%s\r\n", len(msg), msg))
}

func (c *Command) writeArrayBulk(msgs ...string) {
	c.Writer.WriteString(fmt.Sprintf("*%d\r\n", len(msgs)))
	for _, s := range msgs {
		c.Writer.WriteString(fmt.Sprintf("$%d\r\n%s\r\n", len(s), s))
	}
}

func (c *Command) writeNil() {
	c.Writer.WriteString("$-1\r\n")
}
