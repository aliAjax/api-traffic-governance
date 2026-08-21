package repository

import (
	"context"
	"errors"
	"example.com/api-traffic-governance/internal/domain"
	"os"
	"path/filepath"
	"testing"
)

func TestNewEmptyStateInitializesAllMaps(t *testing.T) {
	p := filepath.Join(t.TempDir(), "state.json")
	if err := os.WriteFile(p, []byte(`{"services":null,"upstreams":null,"routes":null,"quotas":null,"experiments":null}`), 0600); err != nil {
		t.Fatal(err)
	}
	st, err := New(p)
	if err != nil {
		t.Fatal(err)
	}
	if err := st.CreateService(context.Background(), domain.Service{ID: "svc", TenantID: "t", Name: "n"}); err != nil {
		t.Fatal(err)
	}
	if _, err := st.GetService(context.Background(), "missing"); !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("missing error = %v", err)
	}
}
