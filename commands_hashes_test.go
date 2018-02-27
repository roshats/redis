package redis

import "testing"

func TestHgetCommand(t *testing.T) {
	t.Parallel()

	t.Run("Key doesn't provided", func(t *testing.T) {
		t.Parallel()

		s := NewMockStorage()
		result := hgetCommand(s, nil)
		_, ok := result.(*ErrorResult)
		Require(t, ok, "Should return error")
	})

	t.Run("Hash key doesn't provided", func(t *testing.T) {
		t.Parallel()

		s := NewMockStorage()
		result := hgetCommand(s, []string{"key"})
		_, ok := result.(*ErrorResult)
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
		str, ok := result.(StringResult)
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
		str, ok := result.(StringResult)
		Require(t, ok, "Returns string")
		Assert(t, str == "", "Returns stored value")
	})
}

func TestHgetallCommand(t *testing.T) {
	t.Parallel()

	t.Run("Key doesn't provided", func(t *testing.T) {
		t.Parallel()

		s := NewMockStorage()
		result := hgetallCommand(s, nil)
		_, ok := result.(*ErrorResult)
		Require(t, ok, "Should return error")
	})

	t.Run("Stored not a hash value", func(t *testing.T) {
		t.Parallel()

		s := NewMockStorage()
		s.Values["key"] = []string{}

		result := hgetallCommand(s, []string{"key"})
		Assert(t, result == wrongValueType, "Should return error")
	})

	t.Run("Key doesn't exist", func(t *testing.T) {
		t.Parallel()

		s := NewMockStorage()
		result := hgetallCommand(s, []string{"key"})
		list, ok := result.(ArrayResult)
		Require(t, ok, "Returns array result")
		Assert(t, len(list) == 0, "Returns empty list")
	})

	t.Run("Returns keys and values", func(t *testing.T) {
		t.Parallel()

		s := NewMockStorage()
		storedDict := map[string]string{"key1": "val1", "key2": "val2"}
		s.Values["key"] = storedDict

		result := hgetallCommand(s, []string{"key"})
		list, ok := result.(ArrayResult)
		Require(t, ok, "Returns array result")
		Require(t, len(list)%2 == 0, "Result has even number of elements")

		for i := 0; i < len(list); i += 2 {
			Assert(t, storedDict[list[i].String()] == list[i+1].String(), "Returns stored value")
		}
	})
}

func TestHsetCommand(t *testing.T) {
	t.Parallel()

	t.Run("Key doesn't provided", func(t *testing.T) {
		t.Parallel()

		s := NewMockStorage()
		result := hsetCommand(s, nil)
		_, ok := result.(*ErrorResult)
		Require(t, ok, "Should return error")
	})

	t.Run("Hash key doesn't provided", func(t *testing.T) {
		t.Parallel()

		s := NewMockStorage()
		result := hsetCommand(s, []string{"key"})
		_, ok := result.(*ErrorResult)
		Require(t, ok, "Should return error")
	})

	t.Run("New value doesn't provided", func(t *testing.T) {
		t.Parallel()

		s := NewMockStorage()
		result := hsetCommand(s, []string{"key", "dictKey"})
		_, ok := result.(*ErrorResult)
		Require(t, ok, "Should return error")
	})

	t.Run("Set value for new key", func(t *testing.T) {
		t.Parallel()

		s := NewMockStorage()
		result := hsetCommand(s, []string{"key", "dictKey", "value"})
		value, ok := result.(IntResult)
		Require(t, ok, "Return int result")
		Assert(t, value == 1, "New hash key")

		Assert(t, s.Values["key"].(map[string]string)["dictKey"] == "value", "Store new value")
	})

	t.Run("Set value for existing key", func(t *testing.T) {
		t.Parallel()

		s := NewMockStorage()
		s.Values["key"] = map[string]string{"dictKey": "value"}

		result := hsetCommand(s, []string{"key", "dictKey", "newVal"})
		value, ok := result.(IntResult)
		Require(t, ok, "Return int result")
		Assert(t, value == 0, "New hash key")

		Assert(t, s.Values["key"].(map[string]string)["dictKey"] == "newVal", "Store new value")
	})
}

func TestHdelCommand(t *testing.T) {
	t.Parallel()

	t.Run("Key doesn't provided", func(t *testing.T) {
		t.Parallel()

		s := NewMockStorage()
		result := hdelCommand(s, nil)
		_, ok := result.(*ErrorResult)
		Require(t, ok, "Should return error")
	})

	t.Run("Hash key doesn't provided", func(t *testing.T) {
		t.Parallel()

		s := NewMockStorage()
		result := hdelCommand(s, []string{"key"})
		_, ok := result.(*ErrorResult)
		Require(t, ok, "Should return error")
	})

	t.Run("Key doensn't exist", func(t *testing.T) {
		t.Parallel()

		s := NewMockStorage()
		result := hdelCommand(s, []string{"key", "dictKey"})
		value, ok := result.(IntResult)
		Require(t, ok, "Return int result")
		Assert(t, value == 0, "Return number of deleted keys")
	})

	t.Run("Multiple keys to delete", func(t *testing.T) {
		t.Parallel()

		s := NewMockStorage()
		s.Values["key"] = map[string]string{"k1": "v1", "k2": "v2", "k3": "v3"}

		result := hdelCommand(s, []string{"key", "k1", "k2", "unknown"})
		value, ok := result.(IntResult)
		Require(t, ok, "Return int result")
		Assert(t, value == 2, "Return number of deleted keys")

		Assert(t, len(s.Values["key"].(map[string]string)) == 1, "Remove required keys")
		Assert(t, s.Values["key"].(map[string]string)["k3"] == "v3", "Do not change value not from query")
	})

	t.Run("Delete all keys", func(t *testing.T) {
		t.Parallel()

		s := NewMockStorage()
		s.Values["key"] = map[string]string{"k1": "v1", "k2": "v2", "k3": "v3"}

		result := hdelCommand(s, []string{"key", "k1", "k2", "k3"})
		value, ok := result.(IntResult)
		Require(t, ok, "Return int result")
		Assert(t, value == 3, "Return number of deleted keys")

		_, exists := s.Values["key"]
		Assert(t, !exists, "Remove hash when all keys are deleted")
	})
}
