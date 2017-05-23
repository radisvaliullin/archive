package main

import (
	"fmt"
	"strconv"
	"strings"
	"test_task_11/mcache/mcache"
)

//
func parseCommand(cstr string) (*mcache.Command, error) {

	cattrs := strings.Fields(cstr)

	c := &mcache.Command{}

	if len(cattrs) < 1 {
		return nil, fmt.Errorf("empty command")
	}

	cname := strings.ToLower(cattrs[0])
	switch cname {
	case "set":
		// set format: command name, type, key, ttl, value/values/key-values

		// command name
		c.Name = cname
		if len(cattrs) < 4 {
			return nil, fmt.Errorf("set command must have key value ttl arguments")
		}

		// command type, key
		ctype := strings.ToLower(cattrs[1])
		c.Type, c.Key = &ctype, mcache.PStr(cattrs[2])

		// ttl
		ttl, err := strconv.ParseInt(cattrs[3], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("set command, ttl attr must be number")
		}
		c.TTL = &ttl

		// values
		switch ctype {
		case "str":
			if len(cattrs) < 5 {
				return nil, fmt.Errorf("set string have not value argument")
			}
			c.Str = mcache.PStr(cattrs[4])

		case "slice":
			if len(cattrs) < 5 {
				return nil, fmt.Errorf("set slice have not value arguments")
			}
			sl := []string{}
			for _, v := range cattrs[4:] {
				sl = append(sl, v)
			}
			c.Slice = sl
		case "map":
			if len(cattrs) < 6 {
				return nil, fmt.Errorf("set map must have minimum one key val of map")
			}
			m := map[string]string{}
			for i := 4; i < len(cattrs); i += 2 {
				if i+1 >= len(cattrs) {
					break
				}
				m[cattrs[i]] = cattrs[i+1]
			}

		default:
			return nil, fmt.Errorf("set unknown type")
		}

	case "get":
		// command name
		c.Name = cname
		if len(cattrs) < 2 {
			return nil, fmt.Errorf("get command must have key")
		}

		// command key, idx_key
		if len(cattrs) == 2 {
			c.Key = mcache.PStr(cattrs[1])
		} else {
			c.Key, c.IdxKey = mcache.PStr(cattrs[1]), mcache.PStr(cattrs[2])
		}

	case "remove":
		// command name
		c.Name = cname
		if len(cattrs) < 2 {
			return nil, fmt.Errorf("remove command must have key")
		}

		// command key
		c.Key = mcache.PStr(cattrs[1])

	case "keys":
		// command name
		c.Name = cname

	default:
		return nil, fmt.Errorf("unknown command")
	}
	return c, nil
}
