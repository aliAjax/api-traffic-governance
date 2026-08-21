package circuit

import (
	"errors"
	"example.com/api-traffic-governance/internal/domain"
	"sync"
	"time"
)

var ErrOpen = errors.New("circuit open")

type Breaker struct {
	mu           sync.Mutex
	state        domain.State
	failures     int
	threshold    int
	openedAt     time.Time
	cooldown     time.Duration
	halfInFlight bool
}

func New(threshold int, cooldown time.Duration) *Breaker {
	if threshold < 1 {
		threshold = 5
	}
	if cooldown <= 0 {
		cooldown = time.Second
	}
	return &Breaker{state: domain.CircuitClosed, threshold: threshold, cooldown: cooldown}
}
func (b *Breaker) Allow() error {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.state == domain.CircuitHalfOpen && b.halfInFlight {
		return ErrOpen
	}
	if b.state == domain.CircuitOpen {
		if time.Since(b.openedAt) < b.cooldown {
			return ErrOpen
		}
		if b.halfInFlight {
			return ErrOpen
		}
		b.state = domain.CircuitHalfOpen
		b.halfInFlight = true
	}
	return nil
}
func (b *Breaker) Success() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.failures = 0
	b.state = domain.CircuitClosed
	b.halfInFlight = false
}
func (b *Breaker) Failure() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.failures++
	b.halfInFlight = false
	if b.failures >= b.threshold {
		b.state = domain.CircuitOpen
		b.openedAt = time.Now()
	}
}
func (b *Breaker) State() domain.State { b.mu.Lock(); defer b.mu.Unlock(); return b.state }
