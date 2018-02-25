package redis

import (
	"github.com/roshats/redis/internal/respgo"
	"strconv"
)

type intResult int // TODO: make int64

func NewIntResult(n int) intResult {
	return intResult(n)
}

func (n intResult) String() string {
	return strconv.Itoa(int(n))
}

func (n intResult) MarshalRESP() []byte {
	return respgo.EncodeInt(int64(n))
}
