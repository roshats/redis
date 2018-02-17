package redis

import (
	"strconv"
	"strings"
	"time"
)

func getCommand(s Storage, query Query) Result {
	if len(query) != 1 {
		return NewErrorResult(generalErrorPrefix, "wrong number of arguments")
	}

	unlock := s.RLock()
	defer unlock()

	key := query[0]
	value, exists := s.Get(key)
	if !exists {
		return nil
	}
	stringValue, ok := value.(string)
	if !ok {
		return NewErrorResult(wrongTypePrefix, "Operation against a key holding the wrong kind of value")
	}
	return NewStringResult(stringValue)
}

func setCommand(s Storage, query Query) Result {
	key, value, expiration, errResult := parseSetCommand(query)
	if errResult != nil {
		return errResult
	}

	unlock := s.Lock()
	defer unlock()

	s.Set(key, value)
	if expiration > 0 {
		t := time.Now().Add(time.Duration(expiration) * time.Second).Unix()
		s.ExpireAt(key, Timestamp(t))
	} else {
		s.RemoveExpiration(key)
	}
	return OKResult
}

func parseSetCommand(queries Query) (key string, value string, expiration int64, result *errorResult) {
	switch len(queries) {
	case 2:
		key = queries[0]
		value = queries[1]
	case 4:
		key = queries[0]
		value = queries[1]
		if strings.ToUpper(queries[2]) != "EX" {
			result = NewErrorResult(generalErrorPrefix, "syntax error")
			return
		}

		var err error
		expiration, err = strconv.ParseInt(queries[3], 10, 64)
		if err != nil || expiration < 0 {
			result = NewErrorResult(generalErrorPrefix, "invalid expire time is set")
			return
		}
	default:
		if len(queries) < 2 {
			result = NewErrorResult(generalErrorPrefix, "wrong number of arguments")
			return
		}
	}
	return
}

func updateCommand(s Storage, query Query) Result {
	if len(query) != 2 {
		return NewErrorResult(generalErrorPrefix, "wrong number of arguments")
	}

	unlock := s.Lock()
	defer unlock()

	s.Set(query[0], query[1])
	return OKResult
}
