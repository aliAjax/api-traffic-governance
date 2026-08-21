package observability

import (
	"errors"
	"example.com/api-traffic-governance/internal/limiter"
	"testing"
	"time"
)

func TestReservationCommitFlowRejectsReplay(t *testing.T) {
	l := limiter.NewLedger()
	_ = l.Reserve("obs-r1", "k", 1, time.Minute)
	_ = l.Commit("obs-r1")
	if err := l.Commit("obs-r1"); !errors.Is(err, limiter.ErrReservationCommitted) {
		t.Fatalf("replay error = %v", err)
	}
}
func TestReservationRollbackFlowRejectsCommitted(t *testing.T) {
	l := limiter.NewLedger()
	_ = l.Reserve("obs-r2", "k", 1, time.Minute)
	_ = l.Commit("obs-r2")
	if err := l.Rollback("obs-r2"); !errors.Is(err, limiter.ErrReservationCommitted) {
		t.Fatalf("rollback error = %v", err)
	}
}
func TestReservationCommitFlowRecordsTimestamp(t *testing.T) {
	l := limiter.NewLedger()
	_ = l.Reserve("obs-r3", "k", 1, time.Minute)
	_ = l.Commit("obs-r3")
	if l.ListCommittedAt("obs-r3").IsZero() {
		t.Fatal("commit timestamp missing")
	}
}
func TestTokenWaitFlowHonorsZeroDeadline(t *testing.T) {
	b := limiter.New(1, 1)
	_ = b.Allow(1)
	started := time.Now()
	if err := b.Wait(1, 0); !errors.Is(err, limiter.ErrWaitExceeded) {
		t.Fatalf("wait error = %v", err)
	}
	if time.Since(started) > 20*time.Millisecond {
		t.Fatal("zero deadline waited too long")
	}
}
