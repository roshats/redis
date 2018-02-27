package redis

import "github.com/roshats/redis/internal/respgo"

var NilResult = NilResultType(false)

type NilResultType bool

func (NilResultType) String() string {
	return "nil"
}

func (NilResultType) MarshalRESP() []byte {
	return respgo.EncodeNull()
}
