package routing

import (
	"context"
	"example.com/api-traffic-governance/internal/domain"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHealthCheckHonorsTimeoutContext(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { time.Sleep(150 * time.Millisecond); w.WriteHeader(200) }))
	defer srv.Close()
	h := HealthChecker{Client: srv.Client(), Timeout: 20 * time.Millisecond}
	started := time.Now()
	ok, _ := h.Check(context.Background(), domain.Upstream{URL: srv.URL})
	if ok {
		t.Fatal("slow upstream reported healthy")
	}
	if elapsed := time.Since(started); elapsed > 100*time.Millisecond {
		t.Fatalf("timeout took %s", elapsed)
	}
}
