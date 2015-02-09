package crawler

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQueue_NewQueue(t *testing.T) {
	q := NewQueue()
	assert.IsType(t, q, new(Queue))
	assert.Equal(t, len(q.queue), 0)
	assert.Equal(t, len(q.pool), 0)
}

func TestQueue_QueuePush(t *testing.T) {
	q := NewQueue()

	links := []string{"link1", "link2", "link3"}

	var wg sync.WaitGroup
	wg.Add(len(links))
	for _, link := range links {
		go func(link string) {
			defer wg.Done()
			q.QueuePush(link)
		}(link)
	}

	wg.Wait()

	assert.Equal(t, len(q.queue), 3)
}

func TestQueue_QueuePush_WithDuplicates(t *testing.T) {
	q := NewQueue()

	links := []string{"link1", "link2", "link3", "link1"}

	var wg sync.WaitGroup
	wg.Add(len(links))
	for _, link := range links {
		go func(link string) {
			defer wg.Done()
			q.QueuePush(link)
		}(link)
	}

	wg.Wait()

	assert.Equal(t, len(q.queue), 3)
}

func TestQueue_PoolPush(t *testing.T) {
	q := NewQueue()

	links := []string{"link1", "link2", "link3"}

	var wg sync.WaitGroup
	wg.Add(len(links))
	for _, link := range links {
		go func(link string) {
			defer wg.Done()
			q.PoolPush(link)
		}(link)
	}

	wg.Wait()

	assert.Equal(t, len(q.pool), 3)
}

func TestQueue_PoolPush_WithDuplicates(t *testing.T) {
	q := NewQueue()

	links := []string{"link1", "link2", "link3", "link1"}

	var wg sync.WaitGroup
	wg.Add(len(links))
	for _, link := range links {
		go func(link string) {
			defer wg.Done()
			q.PoolPush(link)
		}(link)
	}

	wg.Wait()

	assert.Equal(t, len(q.pool), 3)
}

func TestQueue_Len(t *testing.T) {
	q := NewQueue()
	q.QueuePush("link1")
	q.QueuePush("link2")
	q.QueuePush("link3")

	assert.Equal(t, q.Len(), 3)
}

func TestQueue_QueuePop(t *testing.T) {
	q := NewQueue()
	q.QueuePush("link1")
	q.QueuePush("link2")
	q.QueuePush("link3")
	first, err := q.QueuePop()

	assert.Equal(t, first, "link1")
	assert.Equal(t, len(q.pool), 1)
	assert.NoError(t, err)
}

func TestQueue_QueuePop_EmptyQueue(t *testing.T) {
	q := NewQueue()
	_, err := q.QueuePop()

	assert.Error(t, err)
}

func TestQueue_QueuePop_DuplicateItemInQueue(t *testing.T) {
	q := NewQueue()
	q.QueuePush("link1")
	q.QueuePush("link2")

	_, err := q.QueuePop()
	assert.NoError(t, err)

	_, err = q.QueuePop()
	assert.Error(t, err)
}
