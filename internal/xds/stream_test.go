package xds

import (
	"context"
	"testing"
	"time"
)

func TestUnsubscribeWakesWaitingSubscriber(t *testing.T) {
	h := NewHub()
	sub := h.Subscribe("client-1")
	result := make(chan error, 1)
	start := make(chan struct{})
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()
	go func() { <-start; _, err := h.Wait(ctx, sub); result <- err }()
	go func() { <-start; _, err := h.Wait(ctx, sub); result <- err }()
	close(start)
	h.Unsubscribe("client-1")
	for i := 0; i < 2; i++ {
		if err := <-result; err != context.Canceled {
			t.Fatalf("wait error = %v", err)
		}
	}
}
