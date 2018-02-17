package redis

import "testing"

func Assert(t *testing.T, val bool, format string, args ...interface{}) {
	t.Helper()
	if !val {
		t.Errorf(format, args...)
	}
}

func Require(t *testing.T, val bool, format string, args ...interface{}) {
	t.Helper()
	if !val {
		t.Fatalf(format, args...)
	}
}
