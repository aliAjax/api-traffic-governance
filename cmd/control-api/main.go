package main

import (
	"context"
	"errors"
	"example.com/api-traffic-governance/internal/application"
	"example.com/api-traffic-governance/internal/fault"
	"example.com/api-traffic-governance/internal/observability"
	"example.com/api-traffic-governance/internal/repository"
	"example.com/api-traffic-governance/internal/transport/httpapi"
	"example.com/api-traffic-governance/internal/xds"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	log := observability.New()
	st, e := repository.New(env("TRAFFIC_STATE", "./data/state.json"))
	if e != nil {
		log.Error("store", "error", e)
		os.Exit(1)
	}
	svc := application.New(st)
	h := httpapi.New(svc, st, fault.New(), xds.New(env("TRAFFIC_SNAPSHOT_SECRET", "dev-secret")))
	handler := httpapi.RequestID(httpapi.Recover(log)(httpapi.Logger(log)(h.Routes())))
	srv := &http.Server{Addr: env("TRAFFIC_ADDR", ":8085"), Handler: handler, ReadHeaderTimeout: 5 * time.Second, ReadTimeout: 15 * time.Second, WriteTimeout: 15 * time.Second}
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	go func() {
		<-ctx.Done()
		c, cc := context.WithTimeout(context.Background(), 10*time.Second)
		defer cc()
		_ = srv.Shutdown(c)
	}()
	log.Info("listening", "addr", srv.Addr)
	if e := srv.ListenAndServe(); e != nil && !errors.Is(e, http.ErrServerClosed) {
		log.Error("server", "error", e)
		os.Exit(1)
	}
}
func env(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}
