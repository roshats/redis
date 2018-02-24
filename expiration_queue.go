package redis

import (
	"container/heap"
	pq "github.com/roshats/redis/internal/priority_queue"
)

type ExpirationQueue struct {
	queue pq.PriorityQueue
}

func NewExpirationQueue() *ExpirationQueue {
	return &ExpirationQueue{
		queue: make(pq.PriorityQueue, 0),
	}
}

func (q *ExpirationQueue) Push(key string, t Timestamp) {
	item := &pq.Item{Value: key, Priority: int(t)}
	heap.Push(&q.queue, item)
}

func (q *ExpirationQueue) Pop() (string, Timestamp, bool) {
	if q.queue.Len() == 0 {
		return "", 0, false
	}
	item := heap.Pop(&q.queue).(*pq.Item)
	return item.Value, Timestamp(item.Priority), true
}

func (q *ExpirationQueue) Root() (string, Timestamp, bool) {
	if q.queue.Len() == 0 {
		return "", 0, false
	}
	item := q.queue.Root(false).(*pq.Item)
	return item.Value, Timestamp(item.Priority), true
}
