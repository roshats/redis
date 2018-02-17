package redis

import "testing"

func Assert(t *testing.T, val bool, format string, args ...interface{}) {
	if !val {
		t.Errorf(format, args...)
	}
}

func Require(t *testing.T, val bool, format string, args ...interface{}) {
	if !val {
		t.Fatalf(format, args...)
	}
}
