package application

import (
	"errors"
	"example.com/api-traffic-governance/internal/limiter"
	"testing"
	"time"
)

func TestReservationCommitFlowRejectsReplay(t *testing.T) {
	l := limiter.NewLedger()
	_ = l.Reserve("app-r1", "k", 1, time.Minute)
	_ = l.Commit("app-r1")
	if err := l.Commit("app-r1"); !errors.Is(err, limiter.ErrReservationCommitted) {
		t.Fatalf("replay error = %v", err)
	}
}

func TestReservationRollbackFlowRejectsCommitted(t *testing.T) {
	l := limiter.NewLedger()
	_ = l.Reserve("app-r2", "k", 1, time.Minute)
	_ = l.Commit("app-r2")
	if err := l.Rollback("app-r2"); !errors.Is(err, limiter.ErrReservationCommitted) {
		t.Fatalf("rollback error = %v", err)
	}
}

func TestReservationCommitFlowRecordsTimestamp(t *testing.T) {
	l := limiter.NewLedger()
	_ = l.Reserve("app-r3", "k", 1, time.Minute)
	_ = l.Commit("app-r3")
	if l.ListCommittedAt("app-r3").IsZero() {
		t.Fatal("commit timestamp missing")
	}
}

func TestTokenWaitFlowHonorsZeroDeadline(t *testing.T) {
	b := limiter.New(1, 1)
	_ = b.Allow(1)
	if err := b.Wait(1, 0); !errors.Is(err, limiter.ErrWaitExceeded) {
		t.Fatalf("wait error = %v", err)
	}
}
