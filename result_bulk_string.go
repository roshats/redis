package redis

import "github.com/roshats/redis/internal/respgo"

type stringResult string

func NewStringResult(s string) stringResult {
	return stringResult(s)
}

func (s stringResult) String() string {
	return string(s)
}

func (s stringResult) MarshalRESP() []byte {
	return respgo.EncodeBulkString(string(s))
}
