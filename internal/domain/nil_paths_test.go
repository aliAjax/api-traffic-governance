package domain

import "testing"

func TestServiceValidateRequiresName(t *testing.T) {
	if err := (Service{ID: "svc", TenantID: "tenant"}).Validate(); err == nil {
		t.Fatal("service without name was accepted")
	}
}

func TestNilHeaderRuleDoesNotMatch(t *testing.T) {
	if HeaderMatch(nil, map[string]string{"X-Tenant-ID": "t"}) {
		t.Fatal("nil header rule matched a provided header")
	}
}

func TestUpstreamValidateRequiresURL(t *testing.T) {
	if err := (Upstream{ID: "u", ServiceID: "svc", Weight: 1}).Validate(); err == nil {
		t.Fatal("upstream without URL was accepted")
	}
}
