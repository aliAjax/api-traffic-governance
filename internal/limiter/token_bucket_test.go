package limiter

import (
	"testing"
	"time"
)

func TestBucketWaitDoesNotConsumeAfterTimeout(t *testing.T) {
	b := New(1, 1)
	if !b.Allow(1) {
		t.Fatal("initial token unavailable")
	}
	if err := b.Wait(1, 8*time.Millisecond); err == nil {
		t.Fatal("wait unexpectedly succeeded")
	}
	if level := b.Level(); level > 0.1 {
		t.Fatalf("timed out wait restored tokens: %.3f", level)
	}
}
