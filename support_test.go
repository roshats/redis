package redis

import "testing"

func Assert(t *testing.T, val bool) {
	if !val {
		t.Fail()
	}
}

func Require(t *testing.T, val bool) {
	if !val {
		t.FailNow()
	}
}
