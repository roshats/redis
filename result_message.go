package redis

import "github.com/roshats/redis/internal/respgo"

var OKResult = messageResult("OK")

type messageResult string

func (r messageResult) String() string {
	return string(r)
}

func (r messageResult) MarshalRESP() []byte {
	return respgo.EncodeString(string(r))
}
