package queue

import (
	"fmt"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQueue_NewQueue(t *testing.T) {
	q := NewQueue()

	assert.IsType(t, &Queue{}, q)
	assert.Equal(t, 0, len(q.Items))
	assert.Equal(t, "", q.Name)
}

func TestQueue_Len(t *testing.T) {
	q := NewQueue()
	q.Items = append(q.Items, "1", "2", "3")

	assert.Equal(t, 3, q.Len())
}

func TestQueue_Enqueue(t *testing.T) {
	q := NewQueue()

	q.Enqueue("1")
	q.Enqueue("2")
	q.Enqueue("3")

	assert.Equal(t, 3, q.Len())

	q.Enqueue("3")

	assert.Equal(t, 3, q.Len())
}

func TestQueue_Dequeue(t *testing.T) {
	q := NewQueue()
	q.Enqueue("1")
	q.Enqueue("2")
	q.Enqueue("3")

	for i := 1; i < 4; i++ {
		item, err := q.Dequeue()
		assert.Equal(t, item, fmt.Sprintf("%d", i))
		assert.NoError(t, err)
	}

	item, err := q.Dequeue()
	assert.Equal(t, item, "")
	assert.Error(t, err)
}

func TestQueues_NewQueues(t *testing.T) {
	qs := NewQueues()
	assert.Equal(t, 0, len(qs.Queues))
}

func TestQueues_Add(t *testing.T) {
	q := NewQueue()

	qs := NewQueues()
	assert.Equal(t, 0, len(qs.Queues))

	qs.Add(q)
	assert.Equal(t, 1, len(qs.Queues))
}

func TestQueue_Sort(t *testing.T) {
	qs := []*Queue{}

	q := NewQueue()
	q.Name = "google.com"
	qs = append(qs, q)

	q2 := NewQueue()
	q2.Name = "example.com"
	qs = append(qs, q2)

	sort.Sort(QueueList(qs))
	assert.Equal(t, qs, []*Queue{q2, q})
}
