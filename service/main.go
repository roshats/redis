package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/roshats/redis"
	"os"
)

func main() {
	config := NewServerConfig()
	err := updateServerConfig(config, os.Args[1:])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ms := redis.NewMemoryStorage(ctx)
	r := redis.NewServer(ms, redis.CommandsMap, config.Password)
	redis.StartTCPServer(r, config.Port)
}

func updateServerConfig(config *ServerConfig, args []string) error {
	if len(args)%2 != 0 {
		return errors.New("invalid options format")
	}

	for i := 0; i < len(args); i += 2 {
		switch args[i] {
		case "--port":
			config.Port = args[i+1]
		case "--password":
			config.Password = args[i+1]
		default:
			return errors.New(fmt.Sprintf("unknown argument: '%s'", args[i]))
		}
	}
	return nil
}
