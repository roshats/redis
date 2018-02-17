package redis

const (
	generalErrorPrefix = "ERR"
	wrongTypePrefix    = "WRONGTYPE"
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
