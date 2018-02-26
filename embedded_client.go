package redis

import (
	"crypto/subtle"
)

const (
	quitCommand = "quit"
	authCommand = "auth"
)

type EmbeddedClient struct {
	server     ServerInterface
	authorized bool
}

func NewEmbeddedClient(server ServerInterface) *EmbeddedClient {
	return &EmbeddedClient{
		server:     server,
		authorized: len(server.Password()) == 0, // Automatically authorized when password is empty
	}
}

func (c *EmbeddedClient) ProcessCommand(command string, query Query) Result {
	switch command {
	case quitCommand:
		return QuitResult
	case authCommand:
		if len(query) != 1 {
			return wrongNumberOfArgs
		}

		if passwordsMatch([]byte(c.server.Password()), []byte(query[0])) {
			c.authorized = true
			return OKResult
		} else {
			return NewErrorResult(generalErrorPrefix, "wrong password")
		}
	default:
		if !c.authorized {
			return NewErrorResult(generalErrorPrefix, "not authorized")
		}
		return c.server.ProcessCommand(command, query)
	}
}

func passwordsMatch(p1, p2 []byte) bool {
	return subtle.ConstantTimeCompare(p1, p2) == 1
}
