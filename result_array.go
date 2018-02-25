package redis

import (
	"github.com/roshats/redis/internal/respgo"
	"strings"
)

type arrayResult []Result

func (arr arrayResult) String() string {
	stringsArr := make([]string, len(arr))
	for i := range arr {
		stringsArr[i] = arr[i].String()
	}
	return strings.Join(stringsArr, ", ")
}

func (arr arrayResult) MarshalRESP() []byte {
	stringsArr := make([][]byte, len(arr))
	for i := range arr {
		stringsArr[i] = arr[i].MarshalRESP()
	}
	return respgo.EncodeArray(stringsArr)
}

func ArrayResultFromListOfStrings(list []string) Result {
	result := make([]Result, len(list))
	for i := range list {
		result[i] = NewStringResult(list[i])
	}
	return arrayResult(result)
}
