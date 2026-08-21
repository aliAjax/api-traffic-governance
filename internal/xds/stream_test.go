package xds

import (
	"context"
	"testing"
)

func TestUnsubscribeWakesWaitingSubscriber(t *testing.T) {
	h := NewHub()
	sub := h.Subscribe("client-1")
	result := make(chan error, 1)
	go func() { _, err := h.Wait(context.Background(), sub); result <- err }()
	h.Unsubscribe("client-1")
	if err := <-result; err != context.Canceled {
		t.Fatalf("wait error = %v", err)
	}
}
