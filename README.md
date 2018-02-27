# Simple sub-redis implementation

Available commands:
* Strings: `get`, `set`, `update`.
* Keys: `ttl`, `expire`, `del`.
* Lists: `llen`, `lrange`, `ltrim`, `lpush`, `rpush`, `lpop`, `rpop`.
* Hashes: `hget`, `hgetall`, `hset`, `hdel`.
* Connection: `auth`, `quit`.

Run server:

```
$ go run $GOPATH/src/github.com/roshats/redis/service/*.go --password password --port 16379
```

Example usage:

```
$ redis-cli -p 16379
127.0.0.1:16379> set foo bar
(error) ERR not authorized
127.0.0.1:16379> auth password
OK
127.0.0.1:16379> set foo bar ex 10
OK
127.0.0.1:16379> get foo
"bar"
127.0.0.1:16379> ttl foo
(integer) 5
127.0.0.1:16379> get foo
(nil)
127.0.0.1:16379> LPUSH list 1 2 3
(integer) 3
127.0.0.1:16379> RPUSH list 4 5
(integer) 5
127.0.0.1:16379> lrange list 0 -1
1) "3"
2) "2"
3) "1"
4) "4"
5) "5"
127.0.0.1:16379> lpop list
"3"
127.0.0.1:16379> rpop list
"5"
127.0.0.1:16379> hset hash foo bar
(integer) 1
127.0.0.1:16379> hset hash hello world
(integer) 1
127.0.0.1:16379> hgetall hash
1) "foo"
2) "bar"
3) "hello"
4) "world"
127.0.0.1:16379> expire hash 5
(integer) 1
127.0.0.1:16379> ttl hash
(integer) 3
127.0.0.1:16379> ttl hash
(integer) -2
127.0.0.1:16379> set foo bar
OK
127.0.0.1:16379> ttl foo
(integer) -1
127.0.0.1:16379> quit
$
```
