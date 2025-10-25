package command

import (
	"bufio"
	"strings"
	"sync"
)

type Command struct {
	DataMap sync.Map
	Writer  *bufio.Writer
}

type CommandInterface interface {
	HandleCommand(args []string)
}

func NewCommand(w *bufio.Writer) *Command {
	return &Command{Writer: w}
}

type CommandFunc func(c *Command, args []string)

func (c *Command) HandleCommand(cmd []string) {
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

func (c *Command) commands() map[string]CommandFunc {
	return map[string]CommandFunc{
		"echo": func(c *Command, args []string) {
			if len(args) != 1 {
				c.writeError("wrong number of arguments for 'echo'")
				return
			}
			c.writeBulk(args[0])
		},
		"ping": func(c *Command, args []string) {
			c.writeSimple("PONG")
		},
		"set": func(c *Command, args []string) {
			c.handleSetCommand(args)
		},

		"get": func(c *Command, args []string) {
			c.handleGetCommand(args)
		},
	}
}
