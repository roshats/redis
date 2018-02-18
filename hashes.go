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

func hsetCommand(s Storage, query Query) Result {
	if len(query) != 3 {
		return wrongNumberOfArgs
	}

	unlock := s.Lock()
	defer unlock()

	key, dictKey, dictVal := query[0], query[1], query[2]

	dict, errResult := parseDict(s.Get(key))
	if errResult != nil {
		return errResult
	}

	var result int
	if _, existingKey := dict[dictKey]; existingKey {
		result = 0
	} else {
		result = 1
	}

	dict[dictKey] = dictVal
	s.Set(key, dict)
	return NewIntResult(result)
}

func parseDict(value Entry, exists bool) (map[string]string, *errorResult) {
	if !exists {
		return make(map[string]string), nil
	}

	dict, ok := value.(map[string]string)
	if !ok {
		return nil, wrongValueType
	}
	return dict, nil
}
