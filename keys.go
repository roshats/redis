package redis

func ttlCommand(s Storage, query Query) Result {
	if len(query) != 1 {
		return wrongNumberOfArgs
	}

	unlock := s.RLock()
	defer unlock()

	key := query[0]
	expirationTime, exists := s.ExpirationTime(key)
	if !exists {
		if _, valExists := s.Get(key); !valExists {
			return NewIntResult(-2)
		}

		return NewIntResult(-1)
	}

	return NewIntResult(int(expirationTime))
}
