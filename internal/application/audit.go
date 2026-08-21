package application

import (
	"example.com/api-traffic-governance/internal/domain"
	"fmt"
	"sync"
)

type AuditLog struct {
	mu      sync.Mutex
	entries []domain.AuditEntry
}

func NewAuditLog() *AuditLog { return &AuditLog{entries: []domain.AuditEntry{}} }
func (a *AuditLog) Append(e domain.AuditEntry) {
	a.mu.Lock()
	defer a.mu.Unlock()
	if len(a.entries) > 0 {
		e.PreviousHash = a.entries[len(a.entries)-1].Hash
	}
	e.Seal()
	a.entries = append(a.entries, e)
}
func (a *AuditLog) Verify() error {
	a.mu.Lock()
	prev := ""
	for _, e := range a.entries {
		if e.PreviousHash != prev || !e.Verify() {
			return fmt.Errorf("audit chain invalid")
		}
		prev = e.Hash
	}
	a.mu.Unlock()
	return nil
}
func (a *AuditLog) List() []domain.AuditEntry {
	a.mu.Lock()
	defer a.mu.Unlock()
	out := make([]domain.AuditEntry, len(a.entries))
	copy(out, a.entries)
	return out
}
