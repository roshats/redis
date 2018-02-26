package redis

var QuitResult = &quitResultType{OKResult}

type quitResultType struct {
	messageResult
}
