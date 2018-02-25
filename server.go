package redis

type Server struct {
	storage  Storage
	commands map[string]commandFunc
}

func NewServer(storage Storage, commands map[string]commandFunc) *Server {
	return &Server{
		storage:  storage,
		commands: commands,
	}
}

func (s *Server) ProcessCommand(command string, args Query) Result {
	fun, known := s.commands[command]
	if !known {
		return NewErrorResult(generalErrorPrefix, "unknown command")
	}

	return fun(s.storage, args)
}
