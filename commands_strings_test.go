package redis

import (
	"testing"
	"time"
)

func TestGetCommand(t *testing.T) {
	t.Parallel()

	t.Run("Key doesn't provided", func(t *testing.T) {
		t.Parallel()

		s := NewMockStorage()
		result := getCommand(s, nil)
		_, ok := result.(*errorResult)
		Require(t, ok, "Should return error")
	})

	t.Run("Two keys are provided", func(t *testing.T) {
		t.Parallel()

		s := NewMockStorage()
		result := getCommand(s, []string{"key1", "key2"})
		_, ok := result.(*errorResult)
		Require(t, ok, "Should return error")
	})

	t.Run("Provide result when exists", func(t *testing.T) {
		t.Parallel()

		s := NewMockStorage()
		s.Values["key"] = "value"

		result := getCommand(s, []string{"key"})
		str, ok := result.(stringResult)
		Require(t, ok, "Should return string")
		Assert(t, str == "value", "Result should match stored value")
	})

	t.Run("Provide nil result when key doesn't exists", func(t *testing.T) {
		t.Parallel()

		s := NewMockStorage()
		result := getCommand(s, []string{"key"})
		Assert(t, result == NilResult, "Should return nil when associated value found")
	})
}

func TestSetCommand(t *testing.T) {
	t.Parallel()

	t.Run("Value doesn't provided", func(t *testing.T) {
		t.Parallel()

		s := NewMockStorage()
		result := setCommand(s, []string{"key"})
		_, ok := result.(*errorResult)
		Require(t, ok, "Should return error")
		Assert(t, len(s.Values) == 0, "No value should be set when error happened")
		Assert(t, len(s.Expires) == 0, "No expiration should be set when error happened")
	})

	t.Run("Three arguments", func(t *testing.T) {
		t.Parallel()

		s := NewMockStorage()
		result := setCommand(s, []string{"foo", "bar", "baz"})
		_, ok := result.(*errorResult)
		Require(t, ok, "Should return error")
		Assert(t, len(s.Values) == 0, "No value should be set when error happened")
		Assert(t, len(s.Expires) == 0, "No expiration should be set when error happened")
	})

	t.Run("Invalid expiration", func(t *testing.T) {
		t.Parallel()

		s := NewMockStorage()
		result := setCommand(s, []string{"key", "value", "EX", "-1"})
		_, ok := result.(*errorResult)
		Require(t, ok, "Should return error")
		Assert(t, len(s.Values) == 0, "No value should be set when error happened")
		Assert(t, len(s.Expires) == 0, "No expiration should be set when error happened")
	})

	t.Run("Invalid expiration key", func(t *testing.T) {
		t.Parallel()

		s := NewMockStorage()
		result := setCommand(s, []string{"key", "value", "expire", "1"})
		_, ok := result.(*errorResult)
		Require(t, ok, "Should return error")
		Assert(t, len(s.Values) == 0, "No value should be set when error happened")
		Assert(t, len(s.Expires) == 0, "No expiration should be set when error happened")
	})

	t.Run("Set value", func(t *testing.T) {
		t.Parallel()

		s := NewMockStorage()
		result := setCommand(s, []string{"key", "value"})
		Require(t, result == OKResult, "Should return OK on successful insert")

		Assert(t, s.Values["key"] == "value", "Stored value should equal to requested value")
		Assert(t, len(s.Expires) == 0, "Should set no expiration")
	})

	t.Run("Set value with expiration", func(t *testing.T) {
		t.Parallel()

		s := NewMockStorage()
		result := setCommand(s, []string{"key", "value", "EX", "2"})
		Assert(t, result == OKResult, "Should return OK on successful insert")

		Assert(t, s.Values["key"] == "value", "Stored value should equal to requested value")
		Assert(t, AlmostEqual(time.Unix(int64(s.Expires["key"]), 0), TimeAfter(2*time.Second)),
			"Should set expiration")
	})

	t.Run("Set value without expiration when expiration already set", func(t *testing.T) {
		t.Parallel()

		s := NewMockStorage()
		s.Expires["key"] = Timestamp(TimeAfter(time.Minute).Unix())
		result := setCommand(s, []string{"key", "value"})
		Assert(t, result == OKResult, "Should return OK on successful insert")

		Assert(t, s.Values["key"] == "value", "Stored value should equal to requested value")
		Assert(t, len(s.Expires) == 0, "Should remove expiration")
	})

	t.Run("Update expiration", func(t *testing.T) {
		t.Parallel()

		s := NewMockStorage()
		s.Expires["key"] = Timestamp(TimeAfter(time.Minute).Unix())
		result := setCommand(s, []string{"key", "value", "ex", "5"})
		Assert(t, result == OKResult, "Should return OK on successful insert")

		Assert(t, s.Values["key"] == "value", "Stored value should equal to requested value")
		Assert(t, AlmostEqual(time.Unix(int64(s.Expires["key"]), 0), TimeAfter(5*time.Second)),
			"Should update expiration")
	})
}

func TestUpdateCommand(t *testing.T) {
	t.Parallel()

	t.Run("Value doesn't provided", func(t *testing.T) {
		t.Parallel()

		s := NewMockStorage()
		result := updateCommand(s, []string{"key"})
		_, ok := result.(*errorResult)
		Require(t, ok, "Should return error")
		Assert(t, len(s.Values) == 0, "Should not set values")
		Assert(t, len(s.Expires) == 0, "Should not set expiration")
	})

	t.Run("Set value", func(t *testing.T) {
		t.Parallel()

		s := NewMockStorage()
		result := updateCommand(s, []string{"key", "value"})
		Require(t, result == OKResult, "Should return OK on successful update")

		Assert(t, s.Values["key"] == "value", "Stored value should equal to requested value")
		Assert(t, len(s.Expires) == 0, "Should set no expiration")
	})

	t.Run("Doesn't change expiration", func(t *testing.T) {
		t.Parallel()

		s := NewMockStorage()
		result := setCommand(s, []string{"key", "value", "EX", "5"})
		Require(t, result == OKResult, "Should return OK on successful insert")

		result = updateCommand(s, []string{"key", "value"})
		Require(t, result == OKResult, "Should return OK on successful update")

		Assert(t, s.Values["key"] == "value", "Stored value should equal to requested value")
		Assert(t, AlmostEqual(time.Unix(int64(s.Expires["key"]), 0), TimeAfter(5*time.Second)),
			"Should not change expiration")
	})
}
