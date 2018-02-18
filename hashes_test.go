package redis

import "testing"

func TestHgetCommand(t *testing.T) {
	t.Parallel()

	t.Run("Key doesn't provided", func(t *testing.T) {
		t.Parallel()

		s := NewMockStorage()
		result := hgetCommand(s, nil)
		_, ok := result.(*errorResult)
		Require(t, ok, "Should return error")
	})

	t.Run("Hash key doesn't provided", func(t *testing.T) {
		t.Parallel()

		s := NewMockStorage()
		result := hgetCommand(s, []string{"key"})
		_, ok := result.(*errorResult)
		Require(t, ok, "Should return error")
	})

	t.Run("Stored not a hash value", func(t *testing.T) {
		t.Parallel()

		s := NewMockStorage()
		s.Values["key"] = []string{}

		result := hgetCommand(s, []string{"key", "dictKey"})
		Assert(t, result == wrongValueType, "Should return error")
	})

	t.Run("Key doesn't exist", func(t *testing.T) {
		t.Parallel()

		s := NewMockStorage()
		result := hgetCommand(s, []string{"key", "dictKey"})
		Assert(t, result == NilResult, "Returns nil")
	})

	t.Run("Key doesn't exist", func(t *testing.T) {
		t.Parallel()

		s := NewMockStorage()
		result := hgetCommand(s, []string{"key", "dictKey"})
		Assert(t, result == NilResult, "Returns nil")
	})

	t.Run("Returns stored value", func(t *testing.T) {
		t.Parallel()

		s := NewMockStorage()
		s.Values["key"] = map[string]string{"dictKey": "value"}

		result := hgetCommand(s, []string{"key", "dictKey"})
		str, ok := result.(stringResult)
		Require(t, ok, "Returns string")
		Assert(t, str == "value", "Returns stored value")
	})

	t.Run("Empty value", func(t *testing.T) {
		t.Parallel()

		s := NewMockStorage()
		s.Values["key"] = map[string]string{"dictKey": ""}

		result := hgetCommand(s, []string{"key", "anotherDictKey"})
		Assert(t, result == NilResult, "Returns nil")
	})

	t.Run("Empty value", func(t *testing.T) {
		t.Parallel()

		s := NewMockStorage()
		s.Values["key"] = map[string]string{"dictKey": ""}

		result := hgetCommand(s, []string{"key", "dictKey"})
		str, ok := result.(stringResult)
		Require(t, ok, "Returns string")
		Assert(t, str == "", "Returns stored value")
	})
}
