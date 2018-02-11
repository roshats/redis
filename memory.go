package redis

import (
	"sync"
	"time"
)

type MemoryStorage struct {
	// Main storage that holds values
	values map[string]Entry

	// Mapping between keys and expiration time.
	// Expiration time is stored separately from values to save space.
	expires map[string]Timestamp

	mu sync.RWMutex // Do not embed to keep it private
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		values:  make(map[string]Entry),
		expires: make(map[string]Timestamp),
	}
}

func (ms *MemoryStorage) Get(key string) (interface{}, bool) {
	expired := ms.expired(key)
	if expired {
		return nil, false
	}

	val, ok := ms.values[key]
	return val, ok
}

func (ms *MemoryStorage) Set(key string, value interface{}) {
	ms.values[key] = value
}

func (ms *MemoryStorage) Del(key string) {
	delete(ms.values, key)
	ms.RemoveExpiration(key)
}

func (ms *MemoryStorage) ExpireAt(key string, timestamp Timestamp) {
	ms.expires[key] = timestamp
}

func (ms *MemoryStorage) ExpirationTime(key string) (Timestamp, bool) {
	t, ok := ms.expires[key]
	return t, ok
}

func (ms *MemoryStorage) RemoveExpiration(key string) {
	delete(ms.expires, key)
}

func (ms *MemoryStorage) Lock() func() {
	ms.mu.Lock()
	return ms.mu.Unlock
}

func (ms *MemoryStorage) RLock() func() {
	ms.mu.RLock()
	return ms.mu.RUnlock
}

func (ms *MemoryStorage) expired(key string) bool {
	expiration, ok := ms.ExpirationTime(key)
	return ok && time.Unix(int64(expiration), 0).Before(time.Now())
}
