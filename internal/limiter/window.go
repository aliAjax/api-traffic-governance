package limiter

import (
	"sync"
	"time"
)

type Window struct {
	mu       sync.Mutex
	start    time.Time
	count    int
	limit    int
	duration time.Duration
}

func NewWindow(limit int, d time.Duration) *Window {
	if limit < 1 {
		limit = 1
	}
	if d <= 0 {
		d = time.Second
	}
	return &Window{start: time.Now(), limit: limit, duration: d}
}
func (w *Window) Allow() bool {
	w.mu.Lock()
	defer w.mu.Unlock()
	now := time.Now()
	if now.Sub(w.start) >= w.duration {
		w.start = now
		w.count = 0
	}
	if w.count >= w.limit {
		return false
	}
	w.count++
	return true
}
func (w *Window) Remaining() int {
	w.mu.Lock()
	defer w.mu.Unlock()
	if time.Since(w.start) >= w.duration {
		return w.limit
	}
	return w.limit - w.count
}
