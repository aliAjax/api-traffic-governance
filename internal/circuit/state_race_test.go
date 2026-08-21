package circuit

import (
	"example.com/api-traffic-governance/internal/domain"
	"sync"
	"testing"
	"time"
)

func TestStateConcurrentWithFailure(t *testing.T) {
	b := New(100, 0)
	start := make(chan struct{})
	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			for j := 0; j < 1000; j++ {
				b.Failure()
				b.Success()
				_ = b.State()
			}
		}()
	}
	close(start)
	wg.Wait()
}

func TestHalfOpenFailureResetsProbeFailures(t *testing.T) {
	b := New(2, time.Millisecond)
	b.mu.Lock()
	b.state = domain.CircuitHalfOpen
	b.failures = 2
	b.mu.Unlock()
	b.Failure()
	if b.failures != 1 {
		t.Fatalf("failures = %d", b.failures)
	}
}
