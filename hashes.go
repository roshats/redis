package redis

func hgetCommand(s Storage, query Query) Result {
	if len(query) != 2 {
		return wrongNumberOfArgs
	}

	unlock := s.RLock()
	defer unlock()

	key, dictKey := query[0], query[1]
	value, exists := s.Get(key)
	if !exists {
		return NilResult
	}

	dict, ok := value.(map[string]string)
	if !ok {
		return wrongValueType
	}

	dictVal, exists := dict[dictKey]
	if !exists {
		return NilResult
	}
	return NewStringResult(dictVal)
}
