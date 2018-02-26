package redis

import "testing"

func TestEmbeddedClientQuit(t *testing.T) {
	t.Parallel()

	ms := NewMockStorage()
	s := NewServer(ms, nil, "")
	cl := NewEmbeddedClient(s)

	result := cl.ProcessCommand("quit", nil)
	Assert(t, result == QuitResult, "Returns quit result")
}

func TestEmbeddedClientProxyCommand(t *testing.T) {
	t.Parallel()

	ms := NewMockStorage()
	s := NewServer(ms, CommandsMap, "")
	cl := NewEmbeddedClient(s)

	cl.ProcessCommand("set", []string{"hello", "world"})
	result := cl.ProcessCommand("get", []string{"hello"})
	Assert(t, result.String() == "world", "Return stored value")
}

func TestEmbeddedClientAuth(t *testing.T) {
	t.Parallel()

	ms := NewMockStorage()
	s := NewServer(ms, CommandsMap, "password")

	t.Run("Do not process without authorization", func(t *testing.T) {
		cl := NewEmbeddedClient(s)

		result := cl.ProcessCommand("set", []string{"hello", "world"})
		_, ok := result.(*errorResult)
		Assert(t, ok, "Return error")
	})

	t.Run("Wrong password", func(t *testing.T) {
		cl := NewEmbeddedClient(s)

		result := cl.ProcessCommand("auth", []string{"wrong"})
		_, ok := result.(*errorResult)
		Assert(t, ok, "Return error")
	})

	t.Run("Correct password", func(t *testing.T) {
		cl := NewEmbeddedClient(s)

		result := cl.ProcessCommand("auth", []string{"password"})
		Assert(t, result == OKResult, "Return error")

		result = cl.ProcessCommand("set", []string{"hello", "world"})
		Assert(t, result == OKResult, "Process command")
	})

}
