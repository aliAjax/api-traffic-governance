package fault

import (
	"context"
	"example.com/api-traffic-governance/internal/domain"
	"net/http"
	"testing"
	"time"
)

func TestApplyReturnsCanceledContext(t *testing.T) {
	m := New()
	if err := m.Put(domain.Experiment{ID: "e", TenantID: "t", Name: "delay", Scope: "/slow", DelayMs: 10, Status: domain.ExperimentActive}); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	r, _ := http.NewRequest(http.MethodGet, "http://local/slow", nil)
	if _, _, err := m.Apply(ctx, r, "t"); err != context.Canceled {
		t.Fatalf("error = %v", err)
	}
}

func TestApplyReleasesReadLockBeforeDelay(t *testing.T) {
	m := New()
	if err := m.Put(domain.Experiment{ID: "e2", TenantID: "t", Name: "delay", Scope: "/slow", DelayMs: 80, Status: domain.ExperimentActive}); err != nil {
		t.Fatal(err)
	}
	r, _ := http.NewRequest(http.MethodGet, "http://local/slow", nil)
	done := make(chan struct{})
	go func() { _, _, _ = m.Apply(context.Background(), r, "t"); close(done) }()
	time.Sleep(10 * time.Millisecond)
	start := time.Now()
	if err := m.Revoke("e2"); err != nil {
		t.Fatal(err)
	}
	if time.Since(start) > 50*time.Millisecond {
		t.Fatal("revoke waited behind experiment delay")
	}
	<-done
}
