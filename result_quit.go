package redis

var QuitResult = &QuitResultType{OKResult}

type QuitResultType struct {
	MessageResult
}
