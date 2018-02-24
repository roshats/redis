package redis

import (
	"strconv"
	"time"
)

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

func expireCommand(s Storage, query Query) Result {
	if len(query) != 2 {
		return wrongNumberOfArgs
	}

	unlock := s.Lock()
	defer unlock()

	key := query[0]
	expiration, err := strconv.ParseInt(query[1], 10, 64)
	if err != nil || expiration < 0 {
		return NewErrorResult(generalErrorPrefix, "invalid expire time is set")
	}

	_, exists := s.Get(key)
	if !exists {
		return NewIntResult(0)
	}

	expirationTime := time.Now().Add(time.Duration(expiration) * time.Second).Unix()
	s.ExpireAt(key, Timestamp(expirationTime))
	return NewIntResult(1)
}
