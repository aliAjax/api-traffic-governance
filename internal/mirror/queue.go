package mirror

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Message struct {
	ID        string
	TenantID  string
	RouteID   string
	Body      []byte
	CreatedAt time.Time
	ExpiresAt time.Time
}
type Queue struct {
	mu    sync.Mutex
	items []Message
	max   int
}

func NewQueue(max int) *Queue {
	if max < 1 {
		max = 1000
	}
	return &Queue{max: max}
}
func (q *Queue) Push(m Message) error {
	if m.ID == "" || m.TenantID == "" {
		return fmt.Errorf("message identity required")
	}
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.items) >= q.max {
		return fmt.Errorf("mirror queue full")
	}
	q.items = append(q.items, m)
	return nil
}
func (q *Queue) Pop(ctx context.Context) (Message, error) {
	for {
		q.mu.Lock()
		if len(q.items) > 0 {
			m := q.items[0]
			q.items = q.items[1:]
			q.mu.Unlock()
			return m, nil
		}
		q.mu.Unlock()
		select {
		case <-ctx.Done():
			return Message{}, ctx.Err()
		case <-time.After(5 * time.Millisecond):
		}
	}
}
func (q *Queue) Len() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return len(q.items)
}
