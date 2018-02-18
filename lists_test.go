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

func TestLpushCommand(t *testing.T) {
	t.Parallel()

	t.Run("Key doesn't provided", func(t *testing.T) {
		t.Parallel()

		s := NewMockStorage()
		result := lpushCommand(s, nil)
		_, ok := result.(*errorResult)
		Require(t, ok, "Should return error")
	})

	t.Run("No new element is provided", func(t *testing.T) {
		t.Parallel()

		s := NewMockStorage()
		result := lpushCommand(s, []string{"key"})
		_, ok := result.(*errorResult)
		Require(t, ok, "Should return error")
	})

	t.Run("Should create and insert element when no list exists", func(t *testing.T) {
		t.Parallel()

		s := NewMockStorage()
		result := lpushCommand(s, []string{"key", "val1", "val2"})
		value, ok := result.(intResult)
		Require(t, ok, "Return int result")
		Assert(t, value == 2, "Return list length")

		Assert(t, StringsListEqual(s.Values["key"].([]string), []string{"val2", "val1"}),
			"Should store new value")
		Assert(t, len(s.Expires) == 0, "Should not set expiration")
	})

	t.Run("Should insert element when list exists", func(t *testing.T) {
		t.Parallel()

		s := NewMockStorage()
		s.Values["key"] = []string{"value1", "value2"}

		result := lpushCommand(s, []string{"key", "newVal1", "newVal2"})
		value, ok := result.(intResult)
		Require(t, ok, "Return int result")
		Assert(t, value == 4, "Return list length")

		Assert(t, StringsListEqual(s.Values["key"].([]string), []string{"newVal2", "newVal1", "value1", "value2"}),
			"Should store new value")
		Assert(t, len(s.Expires) == 0, "Should not set expiration")
	})
}

func TestRpushCommand(t *testing.T) {
	t.Parallel()

	t.Run("Key doesn't provided", func(t *testing.T) {
		t.Parallel()

		s := NewMockStorage()
		result := rpushCommand(s, nil)
		_, ok := result.(*errorResult)
		Require(t, ok, "Should return error")
	})

	t.Run("No new element is provided", func(t *testing.T) {
		t.Parallel()

		s := NewMockStorage()
		result := rpushCommand(s, []string{"key"})
		_, ok := result.(*errorResult)
		Require(t, ok, "Should return error")
	})

	t.Run("Should create and insert element when no list exists", func(t *testing.T) {
		t.Parallel()

		s := NewMockStorage()
		result := rpushCommand(s, []string{"key", "val1", "val2"})
		value, ok := result.(intResult)
		Require(t, ok, "Return int result")
		Assert(t, value == 2, "Return list length")

		Assert(t, StringsListEqual(s.Values["key"].([]string), []string{"val1", "val2"}),
			"Should store new value")
		Assert(t, len(s.Expires) == 0, "Should not set expiration")
	})

	t.Run("Should insert element when list exists", func(t *testing.T) {
		t.Parallel()

		s := NewMockStorage()
		s.Values["key"] = []string{"value1", "value2"}

		result := rpushCommand(s, []string{"key", "newVal1", "newVal2"})
		value, ok := result.(intResult)
		Require(t, ok, "Return int result")
		Assert(t, value == 4, "Return list length")

		Assert(t, StringsListEqual(s.Values["key"].([]string), []string{"value1", "value2", "newVal1", "newVal2"}),
			"Should store new value")
		Assert(t, len(s.Expires) == 0, "Should not set expiration")
	})
}

func TestLpopCommand(t *testing.T) {
	t.Parallel()

	t.Run("Key doesn't provided", func(t *testing.T) {
		t.Parallel()

		s := NewMockStorage()
		result := lpopCommand(s, nil)
		_, ok := result.(*errorResult)
		Require(t, ok, "Should return error")
	})

	t.Run("Empty list", func(t *testing.T) {
		t.Parallel()

		s := NewMockStorage()
		result := lpopCommand(s, []string{"key"})
		Assert(t, result == NilResult, "Should return nil")
	})

	t.Run("List exists", func(t *testing.T) {
		t.Parallel()

		s := NewMockStorage()
		s.Values["key"] = []string{"val1", "val2", "val3"}

		result := lpopCommand(s, []string{"key"})
		str, ok := result.(stringResult)
		Require(t, ok, "Should return string")
		Assert(t, str == "val1", "Result should match stored value")

		Assert(t, StringsListEqual(s.Values["key"].([]string), []string{"val2", "val3"}),
			"Should update stored value")
	})
}

func TestRpopCommand(t *testing.T) {
	t.Parallel()

	t.Run("Key doesn't provided", func(t *testing.T) {
		t.Parallel()

		s := NewMockStorage()
		result := rpopCommand(s, nil)
		_, ok := result.(*errorResult)
		Require(t, ok, "Should return error")
	})

	t.Run("Empty list", func(t *testing.T) {
		t.Parallel()

		s := NewMockStorage()
		result := rpopCommand(s, []string{"key"})
		Assert(t, result == NilResult, "Should return nil")
	})

	t.Run("List exists", func(t *testing.T) {
		t.Parallel()

		s := NewMockStorage()
		s.Values["key"] = []string{"val1", "val2", "val3"}

		result := rpopCommand(s, []string{"key"})
		str, ok := result.(stringResult)
		Require(t, ok, "Should return string")
		Assert(t, str == "val3", "Result should match stored value")

		Assert(t, StringsListEqual(s.Values["key"].([]string), []string{"val1", "val2"}),
			"Should update stored value")
	})
}
