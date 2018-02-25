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
	expirations      map[string]Timestamp
	expirationsQueue *ExpirationQueue

	mu sync.RWMutex // Do not embed to keep it private
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		values:           make(map[string]Entry),
		expirations:      make(map[string]Timestamp),
		expirationsQueue: NewExpirationQueue(),
	}
}

func (ms *MemoryStorage) Get(key string) (Entry, bool) {
	if t, ok := ms.expirations[key]; ok && timeExpired(t) {
		return nil, false
	}

	val, ok := ms.values[key]
	return val, ok
}

func (ms *MemoryStorage) Set(key string, value Entry) {
	ms.values[key] = value
}

func (ms *MemoryStorage) Del(key string) {
	delete(ms.values, key)
	ms.RemoveExpiration(key)
}

func (ms *MemoryStorage) ExpireAt(key string, timestamp Timestamp) {
	ms.expirationsQueue.Push(key, timestamp)
	ms.expirations[key] = timestamp
}

func (ms *MemoryStorage) ExpirationTime(key string) (Timestamp, bool) {
	if t, ok := ms.expirations[key]; ok && !timeExpired(t) {
		return t, ok
	}
	return 0, false
}

func (ms *MemoryStorage) RemoveExpiration(key string) {
	delete(ms.expirations, key)
}

func (ms *MemoryStorage) Lock() func() {
	ms.mu.Lock()
	return ms.mu.Unlock
}

func (ms *MemoryStorage) RLock() func() {
	ms.mu.RLock()
	return ms.mu.RUnlock
}

func timeExpired(t Timestamp) bool {
	return time.Unix(int64(t), 0).Before(time.Now())
}

func runExpirationCycle(ms *MemoryStorage) {
	// Assume that most of the times there is nothing to expire.
	// Hence use just RLock to check the assumption.
	if !isAnythingToExpire(ms) {
		return
	}

	unlock := ms.Lock()
	defer unlock()

	for {
		k, qt, exists := ms.expirationsQueue.Root()

		// Break expiration cycle when there is no more items in queue or
		// when lowest expiration time is in future.
		if !exists || !timeExpired(qt) {
			break
		}

		ms.expirationsQueue.Pop()
		t, exists := ms.expirations[k]
		if exists && timeExpired(t) {
			ms.Del(k)
		}
	}
}

func isAnythingToExpire(ms *MemoryStorage) bool {
	unlock := ms.RLock()
	defer unlock()

	_, t, exists := ms.expirationsQueue.Root()
	return exists && timeExpired(t)
}
