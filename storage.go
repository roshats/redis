package redis

// Number of seconds elapsed since January 1, 1970 UTC.
type Timestamp int64

type Storage interface {
	Lock() func()
	RLock() func()

	Get(string) (interface{}, bool)
	Set(string, interface{})
	Del(string)

	ExpireAt(string, Timestamp)
	ExpirationTime(string) (Timestamp, bool)
	RemoveExpiration(string)
}
