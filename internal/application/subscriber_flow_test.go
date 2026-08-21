package application

import (
	"context"
	"example.com/api-traffic-governance/internal/xds"
	"testing"
	"time"
)

func TestSubscriberCancellationReturnsThroughApplication(t *testing.T) {
	h := xds.NewHub()
	sub := h.Subscribe("application-client")
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()
	result := make(chan error, 1)
	go func() { _, err := h.Wait(ctx, sub); result <- err }()
	h.Unsubscribe("application-client")
	if err := <-result; err != context.Canceled {
		t.Fatalf("wait error = %v", err)
	}
}
