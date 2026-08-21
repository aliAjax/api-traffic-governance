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

func TestHealthCheckTreatsClientErrorsAsUnhealthy(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusBadRequest) }))
	defer srv.Close()
	ok, _ := (HealthChecker{Client: srv.Client(), Timeout: time.Second}).Check(context.Background(), domain.Upstream{URL: srv.URL})
	if ok {
		t.Fatal("4xx upstream reported healthy")
	}
}
