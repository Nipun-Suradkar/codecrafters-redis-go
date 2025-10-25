package command

import "fmt"

func (c *Command) writeSimple(msg string) {
	c.Writer.WriteString("+" + msg + "\r\n")
	c.Writer.Flush()
}

func (c *Command) writeError(msg string) {
	c.Writer.WriteString("-ERR " + msg + "\r\n")
	c.Writer.Flush()
}

func (c *Command) writeBulk(s string) {
	c.Writer.WriteString(fmt.Sprintf("$%d\r\n%s\r\n", len(s), s))
	c.Writer.Flush()
}

func (c *Command) writeNil() {
	c.Writer.WriteString("$-1\r\n")
	c.Writer.Flush()
}
