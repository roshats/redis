package redis

var (
	OKResult  = singletonResult("OK")
	NilResult = singletonResult("nil")
)

type singletonResult string

func (r singletonResult) String() string {
	return string(r)
}
