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

func hgetallCommand(s Storage, query Query) Result {
	if len(query) != 1 {
		return wrongNumberOfArgs
	}

	unlock := s.RLock()
	defer unlock()

	key := query[0]
	value, exists := s.Get(key)
	if !exists {
		return ArrayResultFromListOfStrings(nil)
	}

	dict, ok := value.(map[string]string)
	if !ok {
		return wrongValueType
	}

	stringsArr := make([]string, 0, 2*len(dict))
	for k, v := range dict {
		stringsArr = append(stringsArr, k, v)
	}
	return ArrayResultFromListOfStrings(stringsArr)
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

func hdelCommand(s Storage, query Query) Result {
	if len(query) < 2 {
		return wrongNumberOfArgs
	}

	unlock := s.Lock()
	defer unlock()

	key, dictKeysToRemove := query[0], query[1:]
	dict, errResult := parseDict(s.Get(key))
	if errResult != nil {
		return errResult
	}

	deleted := 0
	for _, k := range dictKeysToRemove {
		if _, exists := dict[k]; exists {
			deleted += 1
			delete(dict, k)
		}
	}

	if len(dict) > 0 {
		s.Set(key, dict)
	} else {
		s.Del(key)
	}
	return NewIntResult(deleted)
}

func parseDict(value Entry, exists bool) (map[string]string, *ErrorResult) {
	if !exists {
		return make(map[string]string), nil
	}

	dict, ok := value.(map[string]string)
	if !ok {
		return nil, wrongValueType
	}
	return dict, nil
}
