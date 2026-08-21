package mirror

import (
	"testing"
	"time"
)

func TestExpiredDoesNotRewriteInput(t *testing.T) {
	now := time.Now()
	in := []Message{{ID: "live", TenantID: "t", ExpiresAt: now.Add(time.Hour)}, {ID: "old", TenantID: "t", ExpiresAt: now.Add(-time.Hour)}}
	_ = Expired(in, now)
	if len(in) != 2 || in[0].ID != "live" {
		t.Fatalf("input changed: %#v", in)
	}
}

func TestExpiredResultHasIndependentStorage(t *testing.T) {
	now := time.Now()
	in := []Message{{ID: "old", TenantID: "t", ExpiresAt: now.Add(-time.Hour)}}
	out := Expired(in, now)
	out[0].ID = "changed"
	if in[0].ID != "old" {
		t.Fatal("expired result aliases input storage")
	}
}
