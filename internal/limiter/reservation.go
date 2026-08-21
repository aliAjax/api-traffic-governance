package limiter

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

var (
	ErrReservationInvalid   = errors.New("invalid reservation")
	ErrReservationMissing   = errors.New("reservation missing")
	ErrReservationExpired   = errors.New("reservation expired")
	ErrReservationCommitted = errors.New("reservation already committed")
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
		return fmt.Errorf("%w: fields", ErrReservationInvalid)
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	if _, ok := l.items[id]; ok {
		return fmt.Errorf("%w: %s", ErrReservationInvalid, id)
	}
	l.items[id] = Reservation{ID: id, Key: key, Tokens: tokens, ExpiresAt: time.Now().Add(ttl)}
	return nil
}
func (l *Ledger) Commit(id string) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	v, ok := l.items[id]
	if !ok {
		return fmt.Errorf("%w: %s", ErrReservationMissing, id)
	}
	if v.ExpiresAt.Before(time.Now()) {
		return fmt.Errorf("%w: %s", ErrReservationExpired, id)
	}
	v.Committed = true
	l.items[id] = v
	return nil
}
func (l *Ledger) Rollback(id string) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.items, id)
	return nil
}
func (l *Ledger) ListCommittedAt(id string) time.Time { return time.Time{} }
