package main

import (
	"context"
	"github.com/roshats/redis"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ms := redis.NewMemoryStorage(ctx)
	r := redis.NewServer(ms, redis.CommandsMap)
	redis.StartTCPServer(r, "16379")
}
