package redis

import (
	"fmt"
	"testing"
)

func TestLlenCommand(t *testing.T) {
	t.Parallel()

	t.Run("Key doesn't provided", func(t *testing.T) {
		t.Parallel()

		s := NewMockStorage()
		result := llenCommand(s, nil)
		_, ok := result.(*errorResult)
		Require(t, ok, "Should return error")
	})

	t.Run("Two keys are provided", func(t *testing.T) {
		t.Parallel()

		s := NewMockStorage()
		result := llenCommand(s, []string{"key1", "key2"})
		_, ok := result.(*errorResult)
		Require(t, ok, "Should return error")
	})

	t.Run("Return 0 when no key exists", func(t *testing.T) {
		t.Parallel()

		s := NewMockStorage()
		result := llenCommand(s, []string{"key"})
		value, ok := result.(intResult)
		Require(t, ok, "Return int result")
		Assert(t, value == 0, "Return 0 when no key")
	})

	t.Run("Return empty list length", func(t *testing.T) {
		t.Parallel()

		s := NewMockStorage()
		s.Values["key"] = []string{}

		result := llenCommand(s, []string{"key"})
		fmt.Printf("result: %#v\n", result)
		value, ok := result.(intResult)
		Require(t, ok, "Return int result")
		Assert(t, value == 0, "Return 0 for empty list")
	})

	t.Run("Return list with values length", func(t *testing.T) {
		t.Parallel()

		s := NewMockStorage()
		s.Values["key"] = []string{"value1", "value2"}

		result := llenCommand(s, []string{"key"})
		value, ok := result.(intResult)
		Require(t, ok, "Return int result")
		Assert(t, value == 2, "Return list length")
	})
}
