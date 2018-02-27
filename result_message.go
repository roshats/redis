package redis

import "github.com/roshats/redis/internal/respgo"

var OKResult = MessageResult("OK")

type MessageResult string

func (r MessageResult) String() string {
	return string(r)
}

func (r MessageResult) MarshalRESP() []byte {
	return respgo.EncodeString(string(r))
}
