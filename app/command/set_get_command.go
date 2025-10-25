package command

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type SetData struct {
	Value     string
	ValidTill time.Time
}

func (c *Command) handleSetCommand(args []string) {
	if len(args) < 2 {
		c.writeError("wrong number of arguments for 'set'")
	}

	data := &SetData{
		Value: args[1],
	}

	if len(args) == 4 {
		data.ValidTill = getValidTillTime(args)
	}

	c.DataMap.Store(args[0], data)
	c.writeSimple("OK")
}

func (c *Command) handleGetCommand(args []string) {
	if len(args) != 1 {
		c.writeError("wrong number of arguments for 'get'")
		return
	}
	if val, ok := c.DataMap.Load(args[0]); ok {
		if data, ok := val.(*SetData); ok {
			now := time.Now().UTC()
			switch {
			case data.ValidTill.IsZero():
				c.writeBulk(fmt.Sprintf("%v", data.Value))
			case data.ValidTill.After(now):
				c.writeBulk(fmt.Sprintf("%v", data.Value))
			default:
				c.writeNil()
			}
			return
		}
	}
	c.writeNil()
}

func getValidTillTime(args []string) time.Time {
	if strings.EqualFold(args[2], "ex") {
		ttlSeconds, err := strconv.ParseInt(args[3], 10, 64)
		if err != nil {
			return time.Time{}
		}
		duration := time.Duration(ttlSeconds) * time.Second
		return time.Now().UTC().Add(duration)
	}
	if strings.EqualFold(args[2], "px") {
		ttlMilliSeconds, err := strconv.ParseInt(args[3], 10, 64)
		if err != nil {
			return time.Time{}
		}
		duration := time.Duration(ttlMilliSeconds) * time.Millisecond
		return time.Now().UTC().Add(duration)
	}
	return time.Time{}
}
