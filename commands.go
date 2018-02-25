package redis

import "strings"

type commandFunc func(Storage, Query) Result

var CommandsMap = map[string]commandFunc{
	"ttl":    ttlCommand,
	"expire": expireCommand,
	"del":    delCommand,

	"get":    getCommand,
	"set":    setCommand,
	"update": updateCommand,

	"llen":   llenCommand,
	"lrange": lrangeCommand,
	"ltrim":  ltrimCommand,
	"lpush":  lpushCommand,
	"rpush":  rpushCommand,
	"lpop":   lpopCommand,
	"rpop":   rpopCommand,

	"hget":    hgetCommand,
	"hgetall": hgetallCommand,
	"hset":    hsetCommand,
	"hdel":    hdelCommand,
}

func StringToCommand(s string) (string, Query, Result) {
	arr := strings.Fields(s)
	if len(arr) == 0 {
		return "", nil, NewErrorResult(generalErrorPrefix, "Can't parse command")
	}
	return arr[0], arr[1:], nil
}
