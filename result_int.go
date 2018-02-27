package redis

import (
	"github.com/roshats/redis/internal/respgo"
	"strconv"
)

type IntResult int64

func NewIntResult(n int) IntResult {
	return IntResult(n)
}

func (n IntResult) String() string {
	return strconv.Itoa(int(n))
}

func (n IntResult) MarshalRESP() []byte {
	return respgo.EncodeInt(int64(n))
}
