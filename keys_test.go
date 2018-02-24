package redis

import (
	"testing"
	"time"
)

func TestTtlCommand(t *testing.T) {
	t.Parallel()

	t.Run("Key doesn't provided", func(t *testing.T) {
		t.Parallel()

		s := NewMockStorage()
		result := ttlCommand(s, nil)
		_, ok := result.(*errorResult)
		Require(t, ok, "Should return error")
	})

	t.Run("Key doesn't exist", func(t *testing.T) {
		t.Parallel()

		s := NewMockStorage()
		result := ttlCommand(s, []string{"key"})
		val, ok := result.(intResult)
		Require(t, ok, "Should return int")
		Assert(t, val == -2, "Should return -2")
	})

	t.Run("Expiration time is not set", func(t *testing.T) {
		t.Parallel()

		s := NewMockStorage()
		unlock := s.Lock()
		s.Set("key", "value")
		unlock()

		result := ttlCommand(s, []string{"key"})
		val, ok := result.(intResult)
		Require(t, ok, "Should return int")
		Assert(t, val == -1, "Should return -2")
	})

	t.Run("Returns expiration time", func(t *testing.T) {
		t.Parallel()

		s := NewMockStorage()

		expirationTime := TimeAfter(time.Minute).Unix()
		unlock := s.Lock()
		s.Set("key", "value")
		s.ExpireAt("key", Timestamp(expirationTime))
		unlock()

		result := ttlCommand(s, []string{"key"})
		val, ok := result.(intResult)
		Require(t, ok, "Should return int")
		Assert(t, expirationTime == int64(val), "Should return -2")
	})
}
