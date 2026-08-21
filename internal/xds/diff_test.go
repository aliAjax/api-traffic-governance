package xds

import (
	"example.com/api-traffic-governance/internal/domain"
	"testing"
)

func TestCompareDetectsSameLengthQuotaChange(t *testing.T) {
	old := domain.Snapshot{Quotas: []domain.Quota{{ID: "q", RatePerSecond: 1}}}
	next := domain.Snapshot{Quotas: []domain.Quota{{ID: "q", RatePerSecond: 2}}}
	if got := Compare(old, next); got.QuotasAdded == 0 && got.QuotasRemoved == 0 {
		t.Fatal("same-length quota change was ignored")
	}
}
