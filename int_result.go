package redis

import "strconv"

type intResult int

func NewIntResult(n int) intResult {
	return intResult(n)
}

func (n intResult) String() string {
	return strconv.Itoa(int(n))
}
