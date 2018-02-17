package redis

func llenCommand(s Storage, query Query) Result {
	if len(query) != 1 {
		return NewErrorResult(generalErrorPrefix, "wrong number of arguments")
	}

	unlock := s.RLock()
	defer unlock()

	key := query[0]
	value, exists := s.Get(key)
	if !exists {
		// If key does not exist, it is interpreted as an empty list and 0 is returned.
		return NewIntResult(0)
	}

	list, ok := value.([]string)
	if !ok {
		return NewErrorResult(wrongTypePrefix, "Operation against a key holding the wrong kind of value")
	}
	return NewIntResult(len(list))
}
