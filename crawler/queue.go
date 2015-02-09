package crawler

import (
	"errors"
	"sync"
)

type Queue struct {
	queue []string
	pool  []string
	lock  *sync.Mutex
}

func NewQueue() *Queue {
	q := new(Queue)
	q.lock = new(sync.Mutex)
	return q
}

func (q *Queue) QueuePush(link string) {
	q.lock.Lock()
	defer q.lock.Unlock()

	for _, l := range q.queue {
		if l == link {
			return
		}
	}
	q.queue = append(q.queue, link)
}

func (q *Queue) PoolPush(link string) {
	q.lock.Lock()
	defer q.lock.Unlock()

	for _, l := range q.pool {
		if l == link {
			return
		}
	}

	q.pool = append(q.pool, link)
}

func (q *Queue) Len() int {
	q.lock.Lock()
	defer q.lock.Unlock()

	return len(q.queue)
}

func (q *Queue) QueuePop() (string, error) {
	if q.Len() == 0 {
		return "", errors.New("Queue is empty.")
	}

	q.lock.Lock()
	link := q.queue[0]
	for _, l := range q.pool {
		if l == link {
			return "", errors.New("Item is already in pool.")
		}
	}
	q.lock.Unlock()

	q.PoolPush(link)
	return link, nil
}
