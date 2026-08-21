package limiter

import (
	"fmt"
	"sync"
	"time"
)

type Reservation struct {
	ID        string
	Key       string
	Tokens    int
	ExpiresAt time.Time
	Committed bool
}
type Ledger struct {
	mu    sync.Mutex
	items map[string]Reservation
}

func NewLedger() *Ledger { return &Ledger{items: map[string]Reservation{}} }
func (l *Ledger) Reserve(id, key string, tokens int, ttl time.Duration) error {
	if id == "" || key == "" || tokens < 1 {
		return fmt.Errorf("invalid reservation")
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	if _, ok := l.items[id]; ok {
		return fmt.Errorf("reservation exists")
	}
	l.items[id] = Reservation{ID: id, Key: key, Tokens: tokens, ExpiresAt: time.Now().Add(ttl)}
	return nil
}
func (l *Ledger) Commit(id string) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	v, ok := l.items[id]
	if !ok || v.ExpiresAt.Before(time.Now()) {
		return fmt.Errorf("reservation expired")
	}
	if v.Committed {
		return fmt.Errorf("reservation already committed")
	}
	v.Committed = true
	l.items[id] = v
	return nil
}
func (l *Ledger) Rollback(id string) { l.mu.Lock(); defer l.mu.Unlock(); delete(l.items, id) }
