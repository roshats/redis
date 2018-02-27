package redis

import "github.com/roshats/redis/internal/respgo"

type StringResult string

func NewStringResult(s string) StringResult {
	return StringResult(s)
}

func (s StringResult) String() string {
	return string(s)
}

func (s StringResult) MarshalRESP() []byte {
	return respgo.EncodeBulkString(string(s))
}
