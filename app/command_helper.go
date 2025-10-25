package main

import (
	"bufio"
	"fmt"
	"strings"
	"sync"
)

type CommandHelper struct {
	DataMap sync.Map
	Writer  *bufio.Writer
}

type CommandFunc func(c *CommandHelper, args []string)

func NewCommandHelper(w *bufio.Writer) *CommandHelper {
	return &CommandHelper{Writer: w}
}

func (c *CommandHelper) writeSimple(msg string) {
	c.Writer.WriteString("+" + msg + "\r\n")
	c.Writer.Flush()
}

func (c *CommandHelper) writeError(msg string) {
	c.Writer.WriteString("-ERR " + msg + "\r\n")
	c.Writer.Flush()
}

func (c *CommandHelper) writeBulk(s string) {
	c.Writer.WriteString(fmt.Sprintf("$%d\r\n%s\r\n", len(s), s))
	c.Writer.Flush()
}

func (c *CommandHelper) writeNil() {
	c.Writer.WriteString("$-1\r\n")
	c.Writer.Flush()
}

// --- Command registration ---

func (c *CommandHelper) commands() map[string]CommandFunc {
	return map[string]CommandFunc{
		"echo": func(c *CommandHelper, args []string) {
			if len(args) != 1 {
				c.writeError("wrong number of arguments for 'echo'")
				return
			}
			c.writeBulk(args[0])
		},
		"ping": func(c *CommandHelper, args []string) {
			c.writeSimple("PONG")
		},
		"set": func(c *CommandHelper, args []string) {
			if len(args) != 2 {
				c.writeError("wrong number of arguments for 'set'")
				return
			}
			c.DataMap.Store(args[0], args[1])
			c.writeSimple("OK")
		},
		"get": func(c *CommandHelper, args []string) {
			if len(args) != 1 {
				c.writeError("wrong number of arguments for 'get'")
				return
			}
			if val, ok := c.DataMap.Load(args[0]); ok {
				c.writeBulk(fmt.Sprintf("%v", val))
			} else {
				c.writeNil()
			}
		},
	}
}

func (c *CommandHelper) HandleCommand(cmd []string) {
	if len(cmd) == 0 {
		return
	}

	cmdName := strings.ToLower(cmd[0])
	args := cmd[1:]

	if handler, ok := c.commands()[cmdName]; ok {
		handler(c, args)
	} else {
		c.writeError("unknown command '" + cmdName + "'")
	}
}
