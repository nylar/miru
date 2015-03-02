package queue

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQueue_NewQueue(t *testing.T) {
	q := NewQueue()

	assert.IsType(t, &Queue{}, q)
	assert.Equal(t, 0, len(q.Items))
	assert.NotEqual(t, "", q.Name)
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

// func TestQueue_JSON(t *testing.T) {
//  q := NewQueue()
//  q.Enqueue("http://example.com/about/")
//  q.Enqueue("http://example.com/contact/")
//  q.Enqueue("http://example.com/news/")

//  buf := bytes.NewBuffer([]byte{})
//  enc := json.NewEncoder(buf)
//  enc.Encode(q)

//  t.Log(buf.String())
// }

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

// func TestQueues_JSON(t *testing.T) {
//  buf := bytes.NewBuffer([]byte{})
//  enc := json.NewEncoder(buf)

//  q := NewQueue()
//  q.Enqueue("1")
//  q2 := NewQueue()
//  q.Enqueue("2")

//  q.Name = "Queue1"
//  q2.Name = "Queue2"

//  qs := NewQueues()
//  enc.Encode(qs)
//  t.Logf("%s\n", buf.String())

//  qs.Add(q)
//  qs.Add(q2)

//  buf.Reset()
//  enc.Encode(qs)

//  t.Logf("%s\n", buf.String())
// }
