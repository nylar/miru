package miru

import (
	"errors"
	"sync"
)

type Queues struct {
	Queues map[string]*Queue `json:"queues"`
}

func (qs *Queues) Add(q *Queue) {
	qs.Queues[q.Name] = q
}

func NewQueues() *Queues {
	qs := new(Queues)
	qs.Queues = make(map[string]*Queue)
	return qs
}

type QueueList []*Queue

func (ql QueueList) Len() int { return len(ql) }

func (ql QueueList) Swap(i, j int) { ql[i], ql[j] = ql[j], ql[i] }

func (ql QueueList) Less(i, j int) bool { return ql[i].Name < ql[j].Name }

type Queue struct {
	Manager map[string]bool `json:"manager"`
	Items   []string        `json:"items"`
	Name    string          `json:"name"`
	sync.Mutex
}

func NewQueue() *Queue {
	q := new(Queue)
	q.Manager = make(map[string]bool)

	return q
}

func (q *Queue) Enqueue(item string) {
	q.Lock()
	defer q.Unlock()

	if _, ok := q.Manager[item]; !ok {
		q.Manager[item] = true
		q.Items = append(q.Items, item)
	}
}

func (q *Queue) Len() int {
	return len(q.Items)
}

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
