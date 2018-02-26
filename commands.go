package redis

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
