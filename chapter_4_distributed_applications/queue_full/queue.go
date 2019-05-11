package queue_full

import (
	"errors"
	"sync"
)

// ErrMaxQueueSizeReached is returned when the max queue size has been reached.
var ErrMaxQueueSizeReached = errors.New("queue: maximum queue size reached")

// New creates a new queue.
func New(maxSize int) *Queue {
	return &Queue{
		MaxSize: maxSize,
	}
}

// Queue of items.
type Queue struct {
	m       sync.Mutex
	MaxSize int
	Data    []interface{}
}

// Enqueue an item.
func (q *Queue) Enqueue(d interface{}) (err error) {
	q.m.Lock()
	defer q.m.Unlock()
	if len(q.Data) == q.MaxSize {
		return ErrMaxQueueSizeReached
	}
	q.Data = append(q.Data, d)
	return
}

// Dequeue an item.
func (q *Queue) Dequeue() (d interface{}, ok bool) {
	q.m.Lock()
	defer q.m.Unlock()
	if len(q.Data) == 0 {
		return
	}
	ok = true
	d = q.Data[0]
	q.Data = q.Data[1:]
	return
}
