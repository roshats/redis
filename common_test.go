package redis

import "testing"

func TestReverseStrings(t *testing.T) {
	t.Parallel()

	t.Run("nil value", func(t *testing.T) {
		t.Parallel()

		reversed := reverseStrings(nil)
		Assert(t, len(reversed) == 0, "Should return empty list")
	})

	t.Run("List with odd length", func(t *testing.T) {
		t.Parallel()

		reversed := reverseStrings([]string{"a", "b", "c"})
		Assert(t, StringsListEqual(reversed, []string{"c", "b", "a"}),
			"Should reverse list. Got %#v", reversed)
	})

	t.Run("List with even length", func(t *testing.T) {
		t.Parallel()

		reversed := reverseStrings([]string{"a", "b", "c", "d"})
		Assert(t, StringsListEqual(reversed, []string{"d", "c", "b", "a"}),
			"Should reverse list. Got %#v", reversed)
	})
}
