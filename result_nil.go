package redis

import "github.com/roshats/redis/internal/respgo"

var NilResult = nilResultType(false)

type nilResultType bool

func (nilResultType) String() string {
	return "nil"
}

func (nilResultType) MarshalRESP() []byte {
	return respgo.EncodeNull()
}
