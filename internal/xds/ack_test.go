package xds

import (
	"example.com/api-traffic-governance/internal/domain"
	"testing"
)

func TestAckRejectsUnknownVersion(t *testing.T) {
	m := New("secret")
	m.Compile([]domain.Route{{ID: "r1"}}, nil)
	if err := m.Ack("client", "stale-version"); err == nil {
		t.Fatal("stale snapshot was acknowledged")
	}
	if err := m.Ack("client", m.Current().Version); err != nil {
		t.Fatal(err)
	}
}
