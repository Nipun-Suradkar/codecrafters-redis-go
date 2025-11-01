package command

import (
	"bufio"
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/codecrafters-io/redis-starter-go/app/redis_server"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

var once sync.Once

type Command struct {
	Writer *bufio.Writer
}

type CommandInterface interface {
	HandleCommand(args []string)
}

func NewCommand(w *bufio.Writer) *Command {
	once.Do(func() {
		sort.Strings(SupportedCommands)
	})
	return &Command{Writer: w}
}

type CommandFunc func(c *Command, args []string)

func (c *Command) HandleCommand(cmd []string) {
	if len(cmd) == 0 {
		return
	}

	cmdName := strings.ToLower(cmd[0])
	if !isSupportedCommand(cmdName) {
		resp.WriteError(c.Writer, fmt.Sprintf("ERR unknown command '%s'", cmdName))
		return
	}

	args := cmd[1:]

	if handler, ok := c.commands()[cmdName]; ok {
		handler(c, args)
	} else {
		resp.WriteError(c.Writer, "unknown command '"+cmdName+"'")
	}
	c.Writer.Flush()
}

func (c *Command) commands() map[string]CommandFunc {
	return map[string]CommandFunc{
		CmdEcho: func(c *Command, args []string) {
			if len(args) != 1 {
				resp.WriteError(c.Writer, fmt.Sprintf(ErrWrongArgCount, CmdEcho))
				return
			}
			resp.WriteBulkString(c.Writer, args[0])
		},
		CmdPing: func(c *Command, args []string) {
			resp.WriteSimple(c.Writer, "PONG")
		},
		CmdSet: func(c *Command, args []string) {
			c.handleSetCommand(args)
		},
		CmdGet: func(c *Command, args []string) {
			c.handleGetCommand(args)
		},
		CmdConfig: func(c *Command, args []string) {
			if len(args) != 2 {
				resp.WriteError(c.Writer, fmt.Sprintf(ErrWrongArgCount, CmdConfig))
				return
			}
			switch args[1] {
			case ConfigDir:
				resp.WriteArrayBulk(c.Writer, args[1], redis_server.GetRedisServer().DbDir)
			case ConfigDbFile:
				resp.WriteArrayBulk(c.Writer, args[1], redis_server.GetRedisServer().DbDir)
			default:
				resp.WriteNil(c.Writer)
			}
		},
		CmdInfo: func(c *Command, args []string) {
			if len(args) == 1 && strings.EqualFold(args[0], "replication") {
				role := "master"
				if redis_server.GetRedisServer().ReplicaOf != "" {
					role = "slave"
				}
				info := fmt.Sprintf(`role:%s master_replid:%s master_repl_offset:%d`, role, redis_server.GetRedisServer().ReplicationID, redis_server.GetRedisServer().Offset)
				resp.WriteBulkString(c.Writer, info)
			}
		},
		CmdReplicationConfig: func(c *Command, args []string) {
			resp.WriteSimple(c.Writer, "OK")
		},
	}

}

func isSupportedCommand(cmd string) bool {
	cmd = strings.ToLower(cmd)
	i := sort.SearchStrings(SupportedCommands, cmd)
	return i < len(SupportedCommands) && SupportedCommands[i] == cmd
}
