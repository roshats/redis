package redis

import "strings"

type arrayResult []Result

func (arr arrayResult) String() string {
	stringsArr := make([]string, len(arr))
	for i := range arr {
		stringsArr[i] = arr[i].String()
	}
	return strings.Join(stringsArr, ", ")
}

func ArrayResultFromListOfStrings(list []string) Result {
	result := make([]Result, len(list))
	for i := range list {
		result[i] = NewStringResult(list[i])
	}
	return arrayResult(result)
}
