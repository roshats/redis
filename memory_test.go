package redis

import (
	"testing"
	"time"
)

func TestMemoryStorage(t *testing.T) {
	t.Parallel()

	var s interface{} = NewMemoryStorage()
	_, converted := s.(Storage)
	Require(t, converted, "Implements Storage interface")

	t.Run("Set value then get it back", func(t *testing.T) {
		t.Parallel()
		ms := NewMemoryStorage()

		ms.Set("hello", "world")

		val, ok := ms.Get("hello")
		Assert(t, ok, "No value presents in storage")
		Assert(t, val == "world", "Fetched value isn't equal stored value")

		_, temp := ms.ExpirationTime("hello")
		Assert(t, !temp, "No expiration is set")
	})

	t.Run("Multiple keys doesn't conflict", func(t *testing.T) {
		t.Parallel()
		ms := NewMemoryStorage()

		ms.Set("key1", "value1")
		ms.Set("key2", "value2")

		val1, ok := ms.Get("key1")
		Assert(t, ok, "No value at key1")
		Assert(t, val1 == "value1", "Value at key1 isn't equal to stored value")

		val2, ok := ms.Get("key2")
		Assert(t, ok, "No value at key2")
		Assert(t, val2 == "value2", "Value at key2 isn't equal to stored value")
	})

	t.Run("No value for the key", func(t *testing.T) {
		t.Parallel()
		ms := NewMemoryStorage()

		_, ok := ms.Get("unknown")
		Assert(t, !ok, "Value for unknown key is provided")
	})

	t.Run("Do not return expired value", func(t *testing.T) {
		t.Parallel()
		ms := NewMemoryStorage()

		ms.Set("hello", "world")
		ms.ExpireAt("hello", Timestamp(TimeAfter(-time.Second).Unix()))

		_, ok := ms.Get("hello")
		Assert(t, !ok, "Expected to not return value for expired key")
	})

	t.Run("Provide expiration time", func(t *testing.T) {
		t.Parallel()
		ms := NewMemoryStorage()

		ms.Set("hello", "world")
		expected := TimeAfter(2 * time.Minute).Unix()
		ms.ExpireAt("hello", Timestamp(expected))

		actual, ok := ms.ExpirationTime("hello")
		Assert(t, ok, "Expiration time isn't set")
		Assert(t, expected == int64(actual), "Expiration time isn't equal to expected")
	})

	t.Run("No expiration time for expired key", func(t *testing.T) {
		t.Parallel()
		ms := NewMemoryStorage()

		ms.Set("hello", "world")
		ms.ExpireAt("hello", Timestamp(TimeAfter(-time.Second).Unix()))

		_, ok := ms.ExpirationTime("hello")
		Assert(t, !ok, "Should not return expiration time")
	})
}
