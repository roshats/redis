package redis

import "github.com/roshats/redis/internal/respgo"

const (
	generalErrorPrefix = "ERR"
	wrongTypePrefix    = "WRONGTYPE"
)

var (
	wrongNumberOfArgs = NewErrorResult(generalErrorPrefix, "wrong number of arguments")
	wrongValueType    = NewErrorResult(wrongTypePrefix, "Operation against a key holding the wrong kind of value")
)

type errorResult struct {
	prefix, message string
}

func NewErrorResult(prefix, message string) *errorResult {
	return &errorResult{prefix: prefix, message: message}
}

func (r *errorResult) String() string {
	return r.prefix + " " + r.message
}

func (r *errorResult) MarshalRESP() []byte {
	return respgo.EncodeError(r.String())
}
