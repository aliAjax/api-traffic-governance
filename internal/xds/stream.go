package xds

import (
	"context"
	"example.com/api-traffic-governance/internal/domain"
	"sync"
	"time"
)

type Subscriber struct {
	ID   string
	Ch   chan domain.Snapshot
	Done chan struct{}
}
type Hub struct {
	mu   sync.RWMutex
	subs map[string]Subscriber
}

func NewHub() *Hub { return &Hub{subs: map[string]Subscriber{}} }
func (h *Hub) Subscribe(id string) Subscriber {
	h.mu.Lock()
	defer h.mu.Unlock()
	s := Subscriber{ID: id, Ch: make(chan domain.Snapshot, 4), Done: make(chan struct{})}
	h.subs[id] = s
	return s
}
func (h *Hub) Unsubscribe(id string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if _, ok := h.subs[id]; ok {
		delete(h.subs, id)
	}
}
func (h *Hub) Publish(s domain.Snapshot) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	for _, sub := range h.subs {
		select {
		case sub.Ch <- s:
		default:
		}
	}
}
func (h *Hub) Wait(ctx context.Context, sub Subscriber) (domain.Snapshot, error) {
	if sub.ID == "" {
		return domain.Snapshot{}, context.Canceled
	}
	select {
	case s := <-sub.Ch:
		return s, nil
	case <-sub.Done:
		return domain.Snapshot{}, context.Canceled
	case <-time.After(30 * time.Second):
		return domain.Snapshot{}, context.DeadlineExceeded
	case <-ctx.Done():
		return domain.Snapshot{}, ctx.Err()
	}
}
