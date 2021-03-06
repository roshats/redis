package redis

import (
	"testing"
	"time"
)

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

func TimeAfter(duration time.Duration) time.Time {
	return time.Now().Add(duration)
}

func AlmostEqual(t1, t2 time.Time) bool {
	return Within(t1, t2, time.Second)
}

func Within(t1, t2 time.Time, eps time.Duration) bool {
	return t1.Add(-eps).Before(t2) && t1.Add(eps).After(t2)
}

func StringsListEqual(list1, list2 []string) bool {
	if len(list1) != len(list2) {
		return false
	}

	for i := range list1 {
		if list1[i] != list2[i] {
			return false
		}
	}

	return true
}
