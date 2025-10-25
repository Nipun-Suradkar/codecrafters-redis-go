package command

import (
	"bufio"
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/codecrafters-io/redis-starter-go/app/redis_server"
)

var once sync.Once

type Command struct {
	RedisServer *redis_server.RedisServer
	Writer      *bufio.Writer
}

type CommandInterface interface {
	HandleCommand(args []string)
}

func NewCommand(w *bufio.Writer, dataStore *redis_server.RedisServer) *Command {
	once.Do(func() {
		sort.Strings(SupportedCommands)
	})
	return &Command{Writer: w, RedisServer: dataStore}
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
	c.Writer.Flush()
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
				c.writeArrayBulk(args[1], c.RedisServer.DbDir)
			case ConfigDbFile:
				c.writeArrayBulk(args[1], c.RedisServer.DbFilename)
			default:
				c.writeNil()
			}
		},
		CmdInfo: func(c *Command, args []string) {
			if len(args) == 1 && strings.EqualFold(args[0], "replication") {
				role := "master"
				if c.RedisServer.ReplicaOf != "" {
					role = "slave"
				}
				info := fmt.Sprintf(`role:%s master_replid:%s master_repl_offset:%d`, role, c.RedisServer.ReplicationID, c.RedisServer.Offset)
				c.writeBulkString(info)
			}
		},
	}

}

func isSupportedCommand(cmd string) bool {
	cmd = strings.ToLower(cmd)
	i := sort.SearchStrings(SupportedCommands, cmd)
	return i < len(SupportedCommands) && SupportedCommands[i] == cmd
}
