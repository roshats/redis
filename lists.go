package redis

const (
	insertFront = iota
	insertBack
)

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

func lpushCommand(s Storage, query Query) Result {
	return addToList(s, query, insertFront)
}

func rpushCommand(s Storage, query Query) Result {
	return addToList(s, query, insertBack)
}

func addToList(s Storage, query Query, insertMode int) Result {
	if len(query) < 2 {
		return NewErrorResult(generalErrorPrefix, "wrong number of arguments")
	}

	unlock := s.Lock()
	defer unlock()

	key := query[0]

	newValues := query[1:]
	if insertMode == insertFront {
		newValues = reverseStrings(query[1:])
	}

	list, errResult := parseGetList(s.Get(key))
	if errResult != nil {
		return errResult
	}

	head, tail := list, newValues
	if insertMode == insertFront {
		head, tail = tail, head
	}

	newList := append(head, tail...)
	s.Set(key, newList)
	return NewIntResult(len(newList))
}

func lpopCommand(s Storage, query Query) Result {
	return popFromList(s, query, func(list []string) (string, []string) {
		return list[0], list[1:]
	})
}

func rpopCommand(s Storage, query Query) Result {
	return popFromList(s, query, func(list []string) (string, []string) {
		lastIndex := len(list) - 1
		return list[lastIndex], list[:lastIndex]
	})
}

func popFromList(s Storage, query Query, popper func([]string) (string, []string)) Result {
	if len(query) != 1 {
		return NewErrorResult(generalErrorPrefix, "wrong number of arguments")
	}

	unlock := s.Lock()
	defer unlock()

	key := query[0]
	list, errResult := parseGetList(s.Get(key))
	if errResult != nil {
		return errResult
	}
	if len(list) == 0 {
		return NilResult
	}

	result, newList := popper(list)
	s.Set(key, newList)
	return NewStringResult(result)
}

func parseGetList(value Entry, exists bool) ([]string, *errorResult) {
	if !exists {
		return nil, nil
	}

	list, ok := value.([]string)
	if !ok {
		return nil, NewErrorResult(wrongTypePrefix, "Operation against a key holding the wrong kind of value")
	}
	return list, nil
}
