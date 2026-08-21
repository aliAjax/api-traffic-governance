package circuit

import (
	"errors"
	"example.com/api-traffic-governance/internal/domain"
	"testing"
	"time"
)

func TestBreakerHalfOpenAllowsOneProbe(t *testing.T) {
	b := New(1, 5*time.Millisecond)
	if err := b.Allow(); err != nil {
		t.Fatal(err)
	}
	b.Failure()
	if b.State() != domain.CircuitOpen {
		t.Fatalf("state = %s", b.State())
	}
	time.Sleep(8 * time.Millisecond)
	if err := b.Allow(); err != nil {
		t.Fatal(err)
	}
	if err := b.Allow(); !errors.Is(err, ErrOpen) {
		t.Fatalf("second probe error = %v", err)
	}
	b.Success()
	if b.State() != domain.CircuitClosed {
		t.Fatalf("state = %s", b.State())
	}
}
