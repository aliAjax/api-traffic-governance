package limiter

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

var ErrWaitExceeded = errors.New("rate limit wait exceeded")

type Bucket struct {
	mu       sync.Mutex
	tokens   float64
	capacity float64
	rate     float64
	last     time.Time
}

func New(capacity, rate int) *Bucket {
	if capacity < 1 {
		capacity = 1
	}
	if rate < 1 {
		rate = 1
	}
	return &Bucket{tokens: float64(capacity), capacity: float64(capacity), rate: float64(rate), last: time.Now()}
}
func (b *Bucket) Allow(n int) bool {
	if n < 1 {
		n = 1
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	now := time.Now()
	b.tokens += now.Sub(b.last).Seconds() * b.rate
	if b.tokens > b.capacity {
		b.tokens = b.capacity
	}
	b.last = now
	if b.tokens < float64(n) {
		return false
	}
	b.tokens -= float64(n)
	return true
}
func (b *Bucket) Wait(n int, max time.Duration) error {
	if max <= 0 {
		max = 50 * time.Millisecond
	}
	start := time.Now()
	for {
		if b.Allow(n) {
			return nil
		}
		if time.Since(start) >= max {
			return fmt.Errorf("%w after %s", ErrWaitExceeded, max)
		}
		time.Sleep(2 * time.Millisecond)
	}
}
func (b *Bucket) Level() float64 { b.mu.Lock(); defer b.mu.Unlock(); return b.tokens }
