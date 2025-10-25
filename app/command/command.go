package command

import (
	"bufio"
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/codecrafters-io/redis-starter-go/app/data_store"
)

var once sync.Once

type Command struct {
	DataStore *data_store.DataStore
	Writer    *bufio.Writer
}

type CommandInterface interface {
	HandleCommand(args []string)
}

func NewCommand(w *bufio.Writer, dataStore *data_store.DataStore) *Command {
	once.Do(func() {
		sort.Strings(SupportedCommands)
	})
	return &Command{Writer: w, DataStore: dataStore}
}

type CommandFunc func(c *Command, args []string)

func (c *Command) HandleCommand(cmd []string) {
	if len(cmd) == 0 {
		return
	}

	cmdName := strings.ToLower(cmd[0])
	if !isSupportedCommand(cmdName) {
		c.writeError(fmt.Sprintf("ERR unknown command '%s'", cmdName))
		return
	}

	args := cmd[1:]

	if handler, ok := c.commands()[cmdName]; ok {
		handler(c, args)
	} else {
		c.writeError("unknown command '" + cmdName + "'")
	}
}

func (c *Command) commands() map[string]CommandFunc {
	return map[string]CommandFunc{
		CmdEcho: func(c *Command, args []string) {
			if len(args) != 1 {
				c.writeError(fmt.Sprintf(ErrWrongArgCount, CmdEcho))
				return
			}
			c.writeBulkString(args[0])
		},
		CmdPing: func(c *Command, args []string) {
			c.writeSimple("PONG")
		},
		CmdSet: func(c *Command, args []string) {
			c.handleSetCommand(args)
		},
		CmdGet: func(c *Command, args []string) {
			c.handleGetCommand(args)
		},
		CmdConfig: func(c *Command, args []string) {
			if len(args) != 2 {
				c.writeError(fmt.Sprintf(ErrWrongArgCount, CmdConfig))
				return
			}
			switch args[1] {
			case ConfigDir:
				c.writeArrayBulk(args[1], c.DataStore.DbDir)
			case ConfigDbFile:
				c.writeArrayBulk(args[1], c.DataStore.DbFilename)
			default:
				c.writeNil()
			}
		},
		CmdInfo: func(c *Command, args []string) {
			if len(args) == 1 && strings.EqualFold(args[0], "replication") {
				c.writeBulkString("role:master")
			}
		},
	}

}

func isSupportedCommand(cmd string) bool {
	cmd = strings.ToLower(cmd)
	i := sort.SearchStrings(SupportedCommands, cmd)
	return i < len(SupportedCommands) && SupportedCommands[i] == cmd
}
