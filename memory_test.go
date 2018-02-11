package redis

import (
	"testing"
	"time"
)

func TestMemoryStorage(t *testing.T) {
	t.Parallel()

	var s interface{} = NewMemoryStorage()
	_, converted := s.(Storage)
	Require(t, converted)

	t.Run("Set value then get it back", func(t *testing.T) {
		t.Parallel()
		ms := NewMemoryStorage()

		ms.Set("hello", "world")

		val, ok := ms.Get("hello")
		Assert(t, ok)
		Assert(t, val == "world")

		_, temp := ms.ExpirationTime("hello")
		Assert(t, !temp)
	})

	t.Run("Multiple keys doesn't conflict", func(t *testing.T) {
		t.Parallel()
		ms := NewMemoryStorage()

		ms.Set("key1", "value1")
		ms.Set("key2", "value2")

		val1, ok := ms.Get("key1")
		Assert(t, ok)
		Assert(t, val1 == "value1")

		val2, ok := ms.Get("key2")
		Assert(t, ok)
		Assert(t, val2 == "value2")
	})

	t.Run("No value for the key", func(t *testing.T) {
		t.Parallel()
		ms := NewMemoryStorage()

		_, ok := ms.Get("unknown")
		Assert(t, !ok)
	})

	t.Run("Do not return expired value", func(t *testing.T) {
		t.Parallel()
		ms := NewMemoryStorage()

		ms.Set("hello", "world")
		ms.ExpireAt("hello", Timestamp(time.Now().Add(-time.Second).Unix()))

		_, ok := ms.Get("hello")
		Assert(t, !ok)
	})

	t.Run("Provide expiration time", func(t *testing.T) {
		t.Parallel()
		ms := NewMemoryStorage()

		ms.Set("hello", "world")
		expected := time.Now().Add(2 * time.Minute).Unix()
		ms.ExpireAt("hello", Timestamp(expected))

		actual, ok := ms.ExpirationTime("hello")
		Assert(t, ok)
		Assert(t, expected == int64(actual))
	})
}
