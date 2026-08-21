package xds

import (
	"example.com/api-traffic-governance/internal/domain"
	"testing"
)

func TestSnapshotNonceVerifiesCompiledVersion(t *testing.T) {
	m := New("secret")
	s := m.Compile([]domain.Route{{ID: "r1"}}, nil)
	if s.Version == "" || s.Nonce == "" {
		t.Fatal("snapshot credentials are empty")
	}
	if !m.VerifyNonce(s.Nonce) {
		t.Fatal("compiled nonce rejected")
	}
}

func TestSnapshotOwnsInputSlices(t *testing.T) {
	routes := []domain.Route{{ID: "r1"}}
	m := New("secret")
	m.Compile(routes, nil)
	routes[0].ID = "changed"
	if got := m.Current().Routes[0].ID; got != "r1" {
		t.Fatalf("snapshot route changed to %q", got)
	}
}
