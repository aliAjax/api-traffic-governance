package circuit

import "testing"

func TestBudgetAllowsConfiguredRetryBoundary(t *testing.T) {
	if !Budget(Policy{RetryBudget: 2}, 2) {
		t.Fatal("configured retry boundary was rejected")
	}
}
