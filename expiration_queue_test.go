package redis

import "testing"

func TestExpirationQueue(t *testing.T) {
	t.Parallel()

	t.Run("Empty queue", func(t *testing.T) {
		t.Parallel()
		q := NewExpirationQueue()

		_, _, exists := q.Root()
		Assert(t, !exists, "No values in the queue")

		_, _, exists = q.Pop()
		Assert(t, !exists, "No values in the queue")
	})

	t.Run("Root doesn't modify queue", func(t *testing.T) {
		t.Parallel()
		q := NewExpirationQueue()
		q.Push("key", 1)

		for i := 0; i < 2; i++ {
			k, v, exists := q.Root()
			Require(t, exists, "Value in the queue")
			Assert(t, k == "key", "Correct key")
			Assert(t, v == 1, "Correct expiration time")
		}
	})

	t.Run("Pop", func(t *testing.T) {
		t.Parallel()
		q := NewExpirationQueue()

		_, _, exists := q.Root()
		Assert(t, !exists, "No values in the queue")

		q.Push("key", 1)
		_, _, exists = q.Root()
		Assert(t, exists, "Element is in list")

		k, v, exists := q.Pop()
		Require(t, exists, "Value in the queue")
		Assert(t, k == "key", "Correct key")
		Assert(t, v == 1, "Correct expiration time")

		_, _, exists = q.Root()
		Assert(t, !exists, "No values after Pop in the queue")
	})

	t.Run("Root element", func(t *testing.T) {
		t.Parallel()
		q := NewExpirationQueue()

		q.Push("a", 10)
		q.Push("c", 5)

		k1, v1, exists1 := q.Root()
		k, v, exists := q.Pop()
		Assert(t, k == k1, "")
		Assert(t, v == v1, "")
		Assert(t, exists == exists1, "")
	})

	t.Run("Order elements", func(t *testing.T) {
		t.Parallel()
		q := NewExpirationQueue()

		q.Push("a", 10)
		q.Push("c", 5)
		q.Push("b", 20)
		q.Push("d", 7)
		q.Push("e", 15)
		q.Push("f", 2)

		expectedKeys := []string{"f", "c", "d", "a", "e", "b"}
		expectedTimestamps := []Timestamp{2, 5, 7, 10, 15, 20}

		for i := range expectedKeys {
			k, v, exists := q.Pop()
			Require(t, exists, "Value in the queue")

			expectedKey := expectedKeys[i]
			expectedT := expectedTimestamps[i]
			Assert(t, k == expectedKey, "Expected key '%s' but got '%s'", expectedKey, k)
			Assert(t, v == expectedT, "Expected expiration time %d but got %d", expectedT, v)
		}

		_, _, exists := q.Pop()
		Assert(t, !exists, "Queue is empty")
	})
}
