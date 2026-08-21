package limiter

import (
	"errors"
	"testing"
	"time"
)

func TestLedgerRejectsSecondCommit(t *testing.T) {
	l := NewLedger()
	if err := l.Reserve("r1", "tenant/key", 1, time.Minute); err != nil {
		t.Fatal(err)
	}
	if err := l.Commit("r1"); err != nil {
		t.Fatal(err)
	}
	if err := l.Commit("r1"); !errors.Is(err, ErrReservationCommitted) {
		t.Fatal("second commit accepted")
	}
}

func TestLedgerCommitReportsExpiredAndMissing(t *testing.T) {
	l := NewLedger()
	if err := l.Reserve("expired", "tenant/key", 1, -time.Second); err != nil {
		t.Fatal(err)
	}
	if err := l.Commit("expired"); !errors.Is(err, ErrReservationExpired) {
		t.Fatalf("expired error = %v", err)
	}
	if err := l.Commit("missing"); !errors.Is(err, ErrReservationMissing) {
		t.Fatalf("missing error = %v", err)
	}
	if err := l.Commit("expired"); !errors.Is(err, ErrReservationExpired) {
		t.Fatalf("repeated expired error = %v", err)
	}
}

func TestLedgerReserveReportsInvalidSentinel(t *testing.T) {
	if err := NewLedger().Reserve("", "tenant/key", 1, time.Minute); !errors.Is(err, ErrReservationInvalid) {
		t.Fatalf("invalid error = %v", err)
	}
}

func TestLedgerRollbackRejectsCommitted(t *testing.T) {
	l := NewLedger()
	_ = l.Reserve("r3", "k", 1, time.Minute)
	_ = l.Commit("r3")
	if err := l.Rollback("r3"); err == nil {
		t.Fatal("rollback of committed reservation succeeded")
	}
}

func TestLedgerCommitRecordsCommitTime(t *testing.T) {
	l := NewLedger()
	_ = l.Reserve("r2", "k", 1, time.Minute)
	_ = l.Commit("r2")
	if l.items["r2"].CommittedAt.IsZero() {
		t.Fatal("commit time missing")
	}
}
