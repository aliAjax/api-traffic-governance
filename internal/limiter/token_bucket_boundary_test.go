package limiter

import (
	"errors"
	"testing"
	"time"
)

func TestBucketWaitZeroDeadlineReturnsImmediately(t *testing.T) {
	b := New(1, 1)
	_ = b.Allow(1)
	started := time.Now()
	if err := b.Wait(1, 0); !errors.Is(err, ErrWaitExceeded) {
		t.Fatal("zero-deadline wait succeeded")
	}
	if time.Since(started) > 20*time.Millisecond {
		t.Fatal("zero-deadline wait blocked")
	}
}
