package application

import (
	"example.com/api-traffic-governance/internal/domain"
	"testing"
	"time"
)

func TestAuditDetectsAfterFieldTampering(t *testing.T) {
	a := NewAuditLog()
	a.Append(domain.AuditEntry{ID: "1", Action: "publish", Actor: "ops", Resource: "route/a", Before: "draft", After: "published"})
	entries := a.List()
	entries[0].After = "archived"
	b := NewAuditLog()
	b.Append(entries[0])
	entries = b.List()
	entries[0].After = "published"
	if entries[0].Verify() {
		t.Fatal("tampered audit entry still verifies")
	}
}

func TestAuditVerifyReleasesLockAfterFailure(t *testing.T) {
	a := NewAuditLog()
	a.Append(domain.AuditEntry{ID: "2", Action: "publish", Actor: "ops", Resource: "route/b", After: "published"})
	a.entries[0].Hash = "bad"
	if err := a.Verify(); err == nil {
		t.Fatal("invalid chain accepted")
	}
	done := make(chan struct{})
	go func() {
		a.Append(domain.AuditEntry{ID: "3", Action: "publish", Actor: "ops", Resource: "route/c"})
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(100 * time.Millisecond):
		t.Fatal("append blocked after verify failure")
	}
}
