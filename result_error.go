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

type ErrorResult struct {
	prefix, message string
}

func NewErrorResult(prefix, message string) *ErrorResult {
	return &ErrorResult{prefix: prefix, message: message}
}

func (r *ErrorResult) String() string {
	return r.prefix + " " + r.message
}

func (r *ErrorResult) MarshalRESP() []byte {
	return respgo.EncodeError(r.String())
}
