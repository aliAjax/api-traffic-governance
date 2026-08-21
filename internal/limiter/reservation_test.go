package limiter

import (
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
	if err := l.Commit("r1"); err == nil {
		t.Fatal("second commit accepted")
	}
}
