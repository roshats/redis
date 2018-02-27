package main

type ServerConfig struct {
	Port     string
	Password string
}

func NewServerConfig() *ServerConfig {
	return &ServerConfig{
		Port:     "16379",
		Password: "",
	}
}
