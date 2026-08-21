package domain

import (
	"testing"
	"time"
)

func TestAuditSealBindsAfterState(t *testing.T) {
	e := AuditEntry{ID: "1", Action: "publish", Actor: "ops", Resource: "r", Before: "draft", After: "published", CreatedAt: time.Unix(10, 0)}
	e.Seal()
	e.After = "archived"
	if e.Verify() {
		t.Fatal("after-state tampering was accepted")
	}
}

func TestAuditSealRejectsUnsealedEntry(t *testing.T) {
	if (AuditEntry{ID: "empty"}).Verify() {
		t.Fatal("unsealed audit entry verified")
	}
}

func TestAuditSealInitializesTimestampBeforeHash(t *testing.T) {
	e := AuditEntry{ID: "zero-time", Action: "publish", Actor: "ops", Resource: "r"}
	e.Seal()
	if !e.Verify() {
		t.Fatal("freshly sealed entry did not verify")
	}
}

func TestAuditSealBindsCreatedAt(t *testing.T) {
	e := AuditEntry{ID: "2", Action: "publish", Actor: "ops", Resource: "r", CreatedAt: time.Unix(10, 0)}
	e.Seal()
	e.CreatedAt = time.Unix(11, 0)
	if e.Verify() {
		t.Fatal("timestamp tampering was accepted")
	}
}
