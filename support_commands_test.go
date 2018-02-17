package redis

// TODO: Use MemoryStorage
type MockStorage struct {
	Values  map[string]interface{}
	Expires map[string]Timestamp

	Locks, RLocks int
}

func NewMockStorage() *MockStorage {
	return &MockStorage{
		Values:  make(map[string]interface{}),
		Expires: make(map[string]Timestamp),
	}
}

func (s *MockStorage) Lock() func() {
	if s.Locks > 0 {
		panic("Can't obtain write lock twice")
	}

	s.Locks += 1
	return func() {
		s.Locks -= 1
	}
}

func (s *MockStorage) RLock() func() {
	s.RLocks += 1
	return func() {
		s.RLocks -= 1
	}
}

func (s *MockStorage) Get(k string) (v Entry, ok bool) {
	s.assertReadLock()

	v, ok = s.Values[k]
	return
}

func (s *MockStorage) Set(k string, v Entry) {
	s.assertWriteLock()
	s.Values[k] = v
}

func (s *MockStorage) Del(k string) {
	s.assertWriteLock()
	delete(s.Values, k)
	delete(s.Expires, k)
}

func (s *MockStorage) ExpireAt(k string, t Timestamp) {
	s.assertWriteLock()
	s.Expires[k] = t
}

func (s *MockStorage) ExpirationTime(k string) (t Timestamp, ok bool) {
	s.assertReadLock()
	t, ok = s.Expires[k]
	return
}

func (s *MockStorage) RemoveExpiration(k string) {
	s.assertWriteLock()
	delete(s.Expires, k)
}

func (s *MockStorage) assertWriteLock() {
	if s.Locks > 0 {
		return
	}
	panic("Can't write without lock")
}

func (s *MockStorage) assertReadLock() {
	if s.RLocks > 0 || s.Locks > 0 {
		return
	}
	panic("Can't read without lock")
}
