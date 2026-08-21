package limiter

import (
	"fmt"
	"sync"
	"time"
)

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
	start := time.Now()
	for {
		if b.Allow(n) {
			return nil
		}
		if time.Since(start) >= max {
			return fmt.Errorf("rate limit wait exceeded")
		}
		time.Sleep(2 * time.Millisecond)
	}
}
func (b *Bucket) Level() float64 { b.mu.Lock(); defer b.mu.Unlock(); return b.tokens }
