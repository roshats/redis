package redis

import "strconv"

const (
	insertFront = iota
	insertBack
)

func llenCommand(s Storage, query Query) Result {
	if len(query) != 1 {
		return wrongNumberOfArgs
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
		return wrongValueType
	}
	return NewIntResult(len(list))
}

func lrangeCommand(s Storage, query Query) Result {
	key, start, end, errResult := parseLrangeQuery(query)
	if errResult != nil {
		return errResult
	}

	unlock := s.RLock()
	defer unlock()

	list, errResult := getSubList(s, key, start, end)
	if errResult != nil {
		return errResult
	}

	return ArrayResultFromListOfStrings(list)
}

func ltrimCommand(s Storage, query Query) Result {
	key, start, end, errResult := parseLrangeQuery(query)
	if errResult != nil {
		return errResult
	}

	unlock := s.Lock()
	defer unlock()

	list, errResult := getSubList(s, key, start, end)
	if errResult != nil {
		return errResult
	}

	if len(list) > 0 {
		s.Set(key, list)
	} else {
		s.Del(key)
	}
	return OKResult
}

func getSubList(s Storage, key string, start, end int) ([]string, *errorResult) {
	list, errResult := parseGetList(s.Get(key))
	if errResult != nil {
		return nil, errResult
	}

	start, end, empty := adjustRangeIndices(len(list), start, end)
	if empty {
		return nil, nil
	}
	return list[start : end+1], nil
}

func parseLrangeQuery(query Query) (key string, start, end int, errResult *errorResult) {
	if len(query) != 3 {
		errResult = wrongNumberOfArgs
		return
	}

	key = query[0]

	var err error
	if start, err = strconv.Atoi(query[1]); err != nil {
		errResult = NewErrorResult(generalErrorPrefix, "value is not an integer or out of range")
		return
	}
	if end, err = strconv.Atoi(query[2]); err != nil {
		errResult = NewErrorResult(generalErrorPrefix, "value is not an integer or out of range")
		return
	}
	return
}

func adjustRangeIndices(length, start, end int) (_ int, _ int, empty bool) {
	// Make start non-negative
	if start < -length {
		start = 0
	} else if start < 0 {
		start = length + start
	}

	if end < 0 {
		end = length + end
	}

	if start > end || start >= length {
		return 0, 0, true
	}
	// Here we know for sure that end is positive since
	// start is non-negative and start <= end.

	if end >= length {
		end = length - 1
	}

	return start, end, false
}

func lpushCommand(s Storage, query Query) Result {
	return addToList(s, query, insertFront)
}

func rpushCommand(s Storage, query Query) Result {
	return addToList(s, query, insertBack)
}

func addToList(s Storage, query Query, insertMode int) Result {
	if len(query) < 2 {
		return wrongNumberOfArgs
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
		return wrongNumberOfArgs
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
		return nil, wrongValueType
	}
	return list, nil
}
