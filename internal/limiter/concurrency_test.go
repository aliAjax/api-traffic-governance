package limiter

import (
	"context"
	"sync"
	"testing"
)

func TestSemaphoreConcurrentReleaseAccounting(t *testing.T) {
	s := NewSemaphore(8)
	var wg sync.WaitGroup
	for i := 0; i < 64; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := s.Acquire(context.Background()); err != nil {
				t.Errorf("acquire: %v", err)
				return
			}
			s.Release()
		}()
	}
	wg.Wait()
	if got := s.Running(); got != 0 {
		t.Fatalf("running count = %d, want 0", got)
	}
}
