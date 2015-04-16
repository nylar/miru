package miru

import (
	"errors"
	"sync"
)

// Queues is a map of queue's
type Queues struct {
	Queues map[string]*Queue `json:"queues"`
}

// Add pushes a new queue onto the queue list
func (qs *Queues) Add(q *Queue) {
	qs.Queues[q.Name] = q
}

// NewQueues return a new queue list
func NewQueues() *Queues {
	qs := new(Queues)
	qs.Queues = make(map[string]*Queue)
	return qs
}

// Queue holds data regarding a queue
type Queue struct {
	Manager map[string]bool `json:"manager"`
	Items   []string        `json:"items"`
	Name    string          `json:"name"`
	Status  string          `json:"status"`
	sync.Mutex
}

// NewQueue creates a new queue and sets its status to active.
func NewQueue() *Queue {
	q := new(Queue)
	q.Manager = make(map[string]bool)
	q.Status = "active"

	return q
}

// Enqueue pushes a new item onto the queue.
func (q *Queue) Enqueue(item string) {
	q.Lock()
	defer q.Unlock()

	if _, ok := q.Manager[item]; !ok {
		q.Manager[item] = true
		q.Items = append(q.Items, item)
	}
}

// Len returns the number of items in the queue.
func (q *Queue) Len() int {
	return len(q.Items)
}

// Dequeue pops an item and returns it
func (q *Queue) Dequeue() (string, error) {
	q.Lock()
	defer q.Unlock()

	if q.Len() == 0 {
		return "", errors.New("Can't dequeue from an empty queue.")
	}

	oldQueue := q.Items
	x := oldQueue[0]
	newQueue := oldQueue[1:len(oldQueue)]
	q.Items = newQueue
	return x, nil
}
