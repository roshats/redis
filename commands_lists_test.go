package redis

import "testing"

func TestLlenCommand(t *testing.T) {
	t.Parallel()

	t.Run("Key doesn't provided", func(t *testing.T) {
		t.Parallel()

		s := NewMockStorage()
		result := llenCommand(s, nil)
		_, ok := result.(*ErrorResult)
		Require(t, ok, "Should return error")
	})

	t.Run("Two keys are provided", func(t *testing.T) {
		t.Parallel()

		s := NewMockStorage()
		result := llenCommand(s, []string{"key1", "key2"})
		_, ok := result.(*ErrorResult)
		Require(t, ok, "Should return error")
	})

	t.Run("Return 0 when no key exists", func(t *testing.T) {
		t.Parallel()

		s := NewMockStorage()
		result := llenCommand(s, []string{"key"})
		value, ok := result.(IntResult)
		Require(t, ok, "Return int result")
		Assert(t, value == 0, "Return 0 when no key")
	})

	t.Run("Return empty list length", func(t *testing.T) {
		t.Parallel()

		s := NewMockStorage()
		s.Values["key"] = []string{}

		result := llenCommand(s, []string{"key"})
		value, ok := result.(IntResult)
		Require(t, ok, "Return int result")
		Assert(t, value == 0, "Return 0 for empty list")
	})

	t.Run("Return list with values length", func(t *testing.T) {
		t.Parallel()

		s := NewMockStorage()
		s.Values["key"] = []string{"value1", "value2"}

		result := llenCommand(s, []string{"key"})
		value, ok := result.(IntResult)
		Require(t, ok, "Return int result")
		Assert(t, value == 2, "Return list length")
	})
}

func TestLrangeCommand(t *testing.T) {
	t.Parallel()

	t.Run("Key doesn't provided", func(t *testing.T) {
		t.Parallel()

		s := NewMockStorage()
		result := lrangeCommand(s, nil)
		_, ok := result.(*ErrorResult)
		Require(t, ok, "Should return error")
	})

	t.Run("Stop argument is not provided", func(t *testing.T) {
		t.Parallel()

		s := NewMockStorage()
		result := lrangeCommand(s, []string{"key", "0"})
		_, ok := result.(*ErrorResult)
		Require(t, ok, "Should return error")
	})

	t.Run("Not int range", func(t *testing.T) {
		t.Parallel()

		s := NewMockStorage()
		result := lrangeCommand(s, []string{"key", "0", "end"})
		_, ok := result.(*ErrorResult)
		Require(t, ok, "Should return error. Got %#v", result)
	})

	t.Run("Key doesn't exist", func(t *testing.T) {
		t.Parallel()

		s := NewMockStorage()
		result := lrangeCommand(s, []string{"key", "0", "-1"})
		list, ok := result.(ArrayResult)
		Require(t, ok, "Returns array result")
		Assert(t, len(list) == 0, "Returns empty list")
	})

	t.Run("Return array", func(t *testing.T) {
		t.Parallel()

		s := NewMockStorage()
		s.Values["key"] = []string{"val1", "val2", "val3"}

		result := lrangeCommand(s, []string{"key", "0", "2"})
		list, ok := result.(ArrayResult)
		Require(t, ok, "Returns array result")

		var stringsList []string
		for _, e := range list {
			stringsList = append(stringsList, e.String())
		}
		Assert(t, StringsListEqual(stringsList, []string{"val1", "val2", "val3"}), "Returns stored list")
	})

	t.Run("Negative start", func(t *testing.T) {
		t.Parallel()

		s := NewMockStorage()
		s.Values["key"] = []string{"val1", "val2", "val3"}

		result := lrangeCommand(s, []string{"key", "-2", "2"})
		list, ok := result.(ArrayResult)
		Require(t, ok, "Returns array result")

		var stringsList []string
		for _, e := range list {
			stringsList = append(stringsList, e.String())
		}
		Assert(t, StringsListEqual(stringsList, []string{"val2", "val3"}), "Returns stored list")
	})

	t.Run("Big negative start", func(t *testing.T) {
		t.Parallel()

		s := NewMockStorage()
		s.Values["key"] = []string{"val1", "val2", "val3"}

		result := lrangeCommand(s, []string{"key", "-100", "2"})
		list, ok := result.(ArrayResult)
		Require(t, ok, "Returns array result")

		var stringsList []string
		for _, e := range list {
			stringsList = append(stringsList, e.String())
		}
		Assert(t, StringsListEqual(stringsList, []string{"val1", "val2", "val3"}), "Returns stored list")
	})

	t.Run("Big start", func(t *testing.T) {
		t.Parallel()

		s := NewMockStorage()
		s.Values["key"] = []string{"val1", "val2", "val3"}

		result := lrangeCommand(s, []string{"key", "100", "-1"})
		list, ok := result.(ArrayResult)
		Require(t, ok, "Returns array result")

		var stringsList []string
		for _, e := range list {
			stringsList = append(stringsList, e.String())
		}
		Assert(t, len(list) == 0, "Returns empty list")
	})

	t.Run("Negative end", func(t *testing.T) {
		t.Parallel()

		s := NewMockStorage()
		s.Values["key"] = []string{"val1", "val2", "val3"}

		result := lrangeCommand(s, []string{"key", "0", "-2"})
		list, ok := result.(ArrayResult)
		Require(t, ok, "Returns array result")

		var stringsList []string
		for _, e := range list {
			stringsList = append(stringsList, e.String())
		}
		Assert(t, StringsListEqual(stringsList, []string{"val1", "val2"}), "Returns stored list")
	})

	t.Run("Big end", func(t *testing.T) {
		t.Parallel()

		s := NewMockStorage()
		s.Values["key"] = []string{"val1", "val2", "val3"}

		result := lrangeCommand(s, []string{"key", "0", "100"})
		list, ok := result.(ArrayResult)
		Require(t, ok, "Returns array result")

		var stringsList []string
		for _, e := range list {
			stringsList = append(stringsList, e.String())
		}
		Assert(t, StringsListEqual(stringsList, []string{"val1", "val2", "val3"}), "Returns stored list")
	})
}

func TestLtrimCommand(t *testing.T) {
	t.Parallel()

	t.Run("Key doesn't provided", func(t *testing.T) {
		t.Parallel()

		s := NewMockStorage()
		result := ltrimCommand(s, nil)
		_, ok := result.(*ErrorResult)
		Require(t, ok, "Should return error")
	})

	t.Run("Stop argument is not provided", func(t *testing.T) {
		t.Parallel()

		s := NewMockStorage()
		result := ltrimCommand(s, []string{"key", "0"})
		_, ok := result.(*ErrorResult)
		Require(t, ok, "Should return error")
	})

	t.Run("Not int range", func(t *testing.T) {
		t.Parallel()

		s := NewMockStorage()
		result := ltrimCommand(s, []string{"key", "0", "end"})
		_, ok := result.(*ErrorResult)
		Require(t, ok, "Should return error. Got %#v", result)
	})

	t.Run("Key doesn't exist", func(t *testing.T) {
		t.Parallel()

		s := NewMockStorage()
		result := ltrimCommand(s, []string{"key", "0", "-1"})
		Assert(t, result == OKResult, "Returns OK")
	})

	t.Run("Return array", func(t *testing.T) {
		t.Parallel()

		s := NewMockStorage()
		s.Values["key"] = []string{"val1", "val2", "val3"}

		result := ltrimCommand(s, []string{"key", "0", "2"})
		Assert(t, result == OKResult, "Returns OK")

		Assert(t, StringsListEqual(s.Values["key"].([]string), []string{"val1", "val2", "val3"}),
			"Updates stored list")
	})

	t.Run("Negative start", func(t *testing.T) {
		t.Parallel()

		s := NewMockStorage()
		s.Values["key"] = []string{"val1", "val2", "val3"}

		result := ltrimCommand(s, []string{"key", "-2", "2"})
		Assert(t, result == OKResult, "Returns OK")

		Assert(t, StringsListEqual(s.Values["key"].([]string), []string{"val2", "val3"}),
			"Updates stored list")
	})

	t.Run("Big negative start", func(t *testing.T) {
		t.Parallel()

		s := NewMockStorage()
		s.Values["key"] = []string{"val1", "val2", "val3"}

		result := ltrimCommand(s, []string{"key", "-100", "2"})
		Assert(t, result == OKResult, "Returns OK")

		Assert(t, StringsListEqual(s.Values["key"].([]string), []string{"val1", "val2", "val3"}),
			"Updates stored list")
	})

	t.Run("Big start", func(t *testing.T) {
		t.Parallel()

		s := NewMockStorage()
		s.Values["key"] = []string{"val1", "val2", "val3"}

		result := ltrimCommand(s, []string{"key", "100", "-1"})
		Assert(t, result == OKResult, "Returns OK")

		_, exists := s.Values["key"]
		Assert(t, !exists, "Return key when result array is empty")
	})

	t.Run("Negative end", func(t *testing.T) {
		t.Parallel()

		s := NewMockStorage()
		s.Values["key"] = []string{"val1", "val2", "val3"}

		result := ltrimCommand(s, []string{"key", "0", "-2"})
		Assert(t, result == OKResult, "Returns OK")

		Assert(t, StringsListEqual(s.Values["key"].([]string), []string{"val1", "val2"}),
			"Updates stored list")
	})

	t.Run("Big end", func(t *testing.T) {
		t.Parallel()

		s := NewMockStorage()
		s.Values["key"] = []string{"val1", "val2", "val3"}

		result := ltrimCommand(s, []string{"key", "0", "100"})
		Assert(t, result == OKResult, "Returns OK")

		Assert(t, StringsListEqual(s.Values["key"].([]string), []string{"val1", "val2", "val3"}),
			"Updates stored list")
	})
}

func TestLpushCommand(t *testing.T) {
	t.Parallel()

	t.Run("Key doesn't provided", func(t *testing.T) {
		t.Parallel()

		s := NewMockStorage()
		result := lpushCommand(s, nil)
		_, ok := result.(*ErrorResult)
		Require(t, ok, "Should return error")
	})

	t.Run("No new element is provided", func(t *testing.T) {
		t.Parallel()

		s := NewMockStorage()
		result := lpushCommand(s, []string{"key"})
		_, ok := result.(*ErrorResult)
		Require(t, ok, "Should return error")
	})

	t.Run("Should create and insert element when no list exists", func(t *testing.T) {
		t.Parallel()

		s := NewMockStorage()
		result := lpushCommand(s, []string{"key", "val1", "val2"})
		value, ok := result.(IntResult)
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
		value, ok := result.(IntResult)
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
		_, ok := result.(*ErrorResult)
		Require(t, ok, "Should return error")
	})

	t.Run("No new element is provided", func(t *testing.T) {
		t.Parallel()

		s := NewMockStorage()
		result := rpushCommand(s, []string{"key"})
		_, ok := result.(*ErrorResult)
		Require(t, ok, "Should return error")
	})

	t.Run("Should create and insert element when no list exists", func(t *testing.T) {
		t.Parallel()

		s := NewMockStorage()
		result := rpushCommand(s, []string{"key", "val1", "val2"})
		value, ok := result.(IntResult)
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
		value, ok := result.(IntResult)
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
		_, ok := result.(*ErrorResult)
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
		str, ok := result.(StringResult)
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
		_, ok := result.(*ErrorResult)
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
		str, ok := result.(StringResult)
		Require(t, ok, "Should return string")
		Assert(t, str == "val3", "Result should match stored value")

		Assert(t, StringsListEqual(s.Values["key"].([]string), []string{"val1", "val2"}),
			"Should update stored value")
	})
}
