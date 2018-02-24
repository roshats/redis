package redis

type stringResult string

func NewStringResult(s string) stringResult {
	return stringResult(s)
}

func (s stringResult) String() string {
	return string(s)
}
