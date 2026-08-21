package main

import (
	"context"
	"example.com/api-traffic-governance/internal/xds"
	"testing"
)

func TestSubscriberCancellationReturnsThroughWorker(t *testing.T) {
	h := xds.NewHub(); sub := h.Subscribe("worker-client")
	ctx, cancel := context.WithCancel(context.Background()); defer cancel()
	result := make(chan error, 1)
	go func() { _, err := h.Wait(ctx, sub); result <- err }()
	h.Unsubscribe("worker-client")
	if err := <-result; err != context.Canceled { t.Fatalf("wait error = %v", err) }
}
