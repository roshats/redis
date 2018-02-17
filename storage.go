package redis

// Number of seconds elapsed since January 1, 1970 UTC.
type Timestamp int64

type Entry interface{}

type Storage interface {
	Lock() func()
	RLock() func()

	Get(string) (Entry, bool)
	Set(string, Entry)
	Del(string)

	ExpireAt(string, Timestamp)
	ExpirationTime(string) (Timestamp, bool)
	RemoveExpiration(string)
}
