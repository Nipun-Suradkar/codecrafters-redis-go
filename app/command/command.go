package command

import (
	"bufio"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/data_store"
)

type Command struct {
	DataStore *data_store.DataStore
	Writer    *bufio.Writer
}

type CommandInterface interface {
	HandleCommand(args []string)
}

func NewCommand(w *bufio.Writer, dataStore *data_store.DataStore) *Command {
	return &Command{Writer: w, DataStore: dataStore}
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

		"config": func(c *Command, args []string) {
			if len(args) != 2 {
				c.writeError("wrong number of arguments for 'config'")
				return
			}
			switch args[1] {
			case "dir":
				c.writeArrayBulk(args[1], c.DataStore.DbDir)
			case "dbfilename":
				c.writeArrayBulk(args[1], c.DataStore.DbFilename)

			default:
				c.writeNil()
			}
		},
	}
}
