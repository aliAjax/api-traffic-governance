package routing

import (
	"context"
	"example.com/api-traffic-governance/internal/domain"
	"net/http"
	"time"
)

type HealthChecker struct {
	Client  *http.Client
	Timeout time.Duration
}

func (h HealthChecker) Check(ctx context.Context, u domain.Upstream) (bool, int) {
	if h.Timeout <= 0 {
		h.Timeout = time.Second
	}
	client := h.Client
	if client == nil {
		client = &http.Client{}
	}
	c, cancel := context.WithTimeout(ctx, h.Timeout)
	defer cancel()
	start := time.Now()
	req, e := http.NewRequestWithContext(c, http.MethodGet, u.URL, nil)
	if e != nil {
		return false, 0
	}
	resp, e := client.Do(req)
	if e != nil {
		return false, int(time.Since(start).Milliseconds())
	}
	_ = resp.Body.Close()
	return resp.StatusCode >= 200 && resp.StatusCode < 500, int(time.Since(start).Milliseconds())
}
