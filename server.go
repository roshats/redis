package redis

type Server struct {
	storage  Storage
	commands map[string]commandFunc
	password string
}

func NewServer(storage Storage, commands map[string]commandFunc, password string) *Server {
	return &Server{
		storage:  storage,
		commands: commands,
		password: password,
	}
}

func (s *Server) ProcessCommand(command string, args Query) Result {
	fun, known := s.commands[command]
	if !known {
		return NewErrorResult(generalErrorPrefix, "unknown command")
	}

	return fun(s.storage, args)
}

func (s *Server) Password() string {
	return s.password
}
