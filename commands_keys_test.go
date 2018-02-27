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
		_, ok := result.(*ErrorResult)
		Require(t, ok, "Should return error")
	})

	t.Run("Key doesn't exist", func(t *testing.T) {
		t.Parallel()

		s := NewMockStorage()
		result := ttlCommand(s, []string{"key"})
		val, ok := result.(IntResult)
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
		val, ok := result.(IntResult)
		Require(t, ok, "Should return int")
		Assert(t, val == -1, "Should return -2")
	})

	t.Run("Returns TTL time", func(t *testing.T) {
		t.Parallel()

		s := NewMockStorage()

		expirationTime := TimeAfter(time.Minute)
		unlock := s.Lock()
		s.Set("key", "value")
		s.ExpireAt("key", Timestamp(expirationTime.Unix()))
		unlock()

		result := ttlCommand(s, []string{"key"})
		val, ok := result.(IntResult)
		Require(t, ok, "Should return int")

		Assert(t, AlmostEqual(time.Now().Add(time.Duration(val)*time.Second), expirationTime),
			"Should return TTL")
	})
}

func TestExpireCommand(t *testing.T) {
	t.Parallel()

	t.Run("Key doesn't provided", func(t *testing.T) {
		t.Parallel()

		s := NewMockStorage()
		result := expireCommand(s, nil)
		_, ok := result.(*ErrorResult)
		Require(t, ok, "Should return error")
	})

	t.Run("Timeout doesn't provided", func(t *testing.T) {
		t.Parallel()

		s := NewMockStorage()
		result := expireCommand(s, []string{"key"})
		_, ok := result.(*ErrorResult)
		Require(t, ok, "Should return error")
	})

	t.Run("Key doesn't exist", func(t *testing.T) {
		t.Parallel()

		s := NewMockStorage()
		result := expireCommand(s, []string{"key", "10"})
		val, ok := result.(IntResult)
		Require(t, ok, "Should return int")
		Assert(t, val == 0, "Should return 0")
	})

	t.Run("Timeout set", func(t *testing.T) {
		t.Parallel()

		s := NewMockStorage()
		s.Values["key"] = "value"

		result := expireCommand(s, []string{"key", "10"})
		val, ok := result.(IntResult)
		Require(t, ok, "Should return int")
		Assert(t, val == 1, "Should return 0")

		Assert(t, AlmostEqual(time.Unix(int64(s.Expires["key"]), 0), TimeAfter(10*time.Second)),
			"Should set expiration")
	})
}

func TestDelCommand(t *testing.T) {
	t.Parallel()

	t.Run("Key doesn't provided", func(t *testing.T) {
		t.Parallel()

		s := NewMockStorage()
		result := delCommand(s, nil)
		_, ok := result.(*ErrorResult)
		Require(t, ok, "Should return error")
	})

	t.Run("Returns keys", func(t *testing.T) {
		t.Parallel()

		s := NewMockStorage()
		s.Values["k1"] = "v1"
		s.Values["k2"] = "v2"
		s.Values["k3"] = "v3"
		s.Expires["k3"] = Timestamp(TimeAfter(time.Minute).Unix())
		s.Values["k4"] = "v4"
		s.Expires["k4"] = Timestamp(TimeAfter(time.Minute).Unix())

		result := delCommand(s, []string{"k1", "k2", "k3", "unknown"})
		val, ok := result.(IntResult)
		Require(t, ok, "Should return int")
		Assert(t, val == 3, "3 keys should be deleted")

		Assert(t, len(s.Values) == 1, "Only one key left")
		_, exists := s.Values["k4"]
		Assert(t, exists, "k4 key left")

		Assert(t, len(s.Expires) == 1, "Only one expiration value left")
		_, exists = s.Expires["k4"]
		Assert(t, exists, "k4 expiration value left")
	})
}
