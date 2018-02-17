package redis

const okString = "OK"

var OKResult okResultType

type okResultType bool

func (okResultType) String() string {
	return okString
}
