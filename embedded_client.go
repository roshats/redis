package redis

const (
	quitCommand = "quit"
)

type EmbeddedClient struct {
	server CommandProcessor
	// TODO: add is authorized
}

func NewEmbeddedClient(server *Server) *EmbeddedClient {
	return &EmbeddedClient{
		server: server,
	}
}

func (c *EmbeddedClient) ProcessCommand(command string, query Query) Result {
	if command == quitCommand {
		return QuitResult
	}
	return c.server.ProcessCommand(command, query)
}
