package circuit

import (
	"net/http"
	"time"
)

type Policy struct {
	FailureThreshold int
	Cooldown         time.Duration
	RetryBudget      int
}

func DefaultPolicy() Policy {
	return Policy{FailureThreshold: 5, Cooldown: time.Second, RetryBudget: 2}
}
func Retryable(status int, err error) bool {
	if err != nil {
		return true
	}
	return status == 408 || status == 425 || status == 429 || status >= 500
}
func Budget(policy Policy, attempt int) bool { return attempt < policy.RetryBudget }
func StatusError(resp *http.Response) bool   { return resp == nil || resp.StatusCode >= 500 }
