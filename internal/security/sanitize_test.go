package security

import "testing"

func TestHeadersDoesNotMutateInput(t *testing.T) {
	in := map[string]string{"Authorization": "Bearer abc", "X-Tenant": "t"}
	out := Headers(in)
	if out["Authorization"] != "[REDACTED]" || in["Authorization"] != "Bearer abc" {
		t.Fatalf("header input changed: in=%v out=%v", in, out)
	}
}
