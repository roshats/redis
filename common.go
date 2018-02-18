package redis

import "fmt"

type Query []string

type Result fmt.Stringer

func reverseStrings(list []string) []string {
	length := len(list)
	result := make([]string, length)

	for i := range list {
		result[i] = list[length-i-1]
	}
	return result
}
