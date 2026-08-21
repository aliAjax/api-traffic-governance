package limiter

import (
	"context"
	"fmt"
	"sync"
)

type Semaphore struct {
	slots   chan struct{}
	mu      sync.Mutex
	running int
}

func NewSemaphore(limit int) *Semaphore {
	if limit < 1 {
		limit = 1
	}
	return &Semaphore{slots: make(chan struct{}, limit)}
}
func (s *Semaphore) Acquire(ctx context.Context) error {
	select {
	case s.slots <- struct{}{}:
		s.mu.Lock()
		s.running++
		s.mu.Unlock()
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
func (s *Semaphore) Release() {
	select {
	case <-s.slots:
		s.mu.Lock()
		if s.running > 0 {
			s.running--
		}
		s.mu.Unlock()
	default:
	}
}
func (s *Semaphore) Running() int { s.mu.Lock(); defer s.mu.Unlock(); return s.running }
func (s *Semaphore) Resize(limit int) error {
	if limit < 1 {
		return fmt.Errorf("limit must be positive")
	}
	return nil
}
