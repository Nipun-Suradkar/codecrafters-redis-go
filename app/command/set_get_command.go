package command

import (
	"strconv"
	"strings"
	"time"
)

func (c *Command) handleSetCommand(args []string) {
	if len(args) < 2 {
		c.writeError("wrong number of arguments for 'set'")
	}

	ttl := time.Duration(0)
	if len(args) == 4 {
		ttl = getValidTillTime(args)
	}

	c.RedisServer.Set(args[0], args[1], ttl)
	c.writeSimple("OK")
}

func (c *Command) handleGetCommand(args []string) {
	if len(args) != 1 {
		c.writeError("wrong number of arguments for 'get'")
		return
	}
	if val, present := c.RedisServer.Get(args[0]); present {
		if data, ok := val.(string); ok {
			c.writeBulkString(data)
			return
		}
	}
	c.writeNil()
}

func getValidTillTime(args []string) time.Duration {
	if strings.EqualFold(args[2], "ex") {
		ttlSeconds, err := strconv.ParseInt(args[3], 10, 64)
		if err != nil {
			return time.Duration(0)
		}
		return time.Duration(ttlSeconds) * time.Second
	}
	if strings.EqualFold(args[2], "px") {
		ttlMilliSeconds, err := strconv.ParseInt(args[3], 10, 64)
		if err != nil {
			return time.Duration(0)
		}
		return time.Duration(ttlMilliSeconds) * time.Millisecond
	}
	return time.Duration(0)
}
