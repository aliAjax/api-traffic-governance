package mirror

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestQueueConcurrentPushRespectsCapacity(t *testing.T) {
	q := NewQueue(4)
	start := make(chan struct{})
	var wg sync.WaitGroup
	for i := 0; i < 32; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			<-start
			_ = q.Push(Message{ID: fmt.Sprint(i), TenantID: "tenant", CreatedAt: time.Now()})
		}(i)
	}
	close(start)
	wg.Wait()
	if got := q.Len(); got > 4 {
		t.Fatalf("queue length = %d, want <= 4", got)
	}
}

func TestQueueLengthConcurrentWithPush(t *testing.T) {
	q := NewQueue(128)
	start := make(chan struct{})
	var wg sync.WaitGroup
	for i := 0; i < 32; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			<-start
			_ = q.Push(Message{ID: fmt.Sprint(i), TenantID: "tenant"})
			_ = q.Len()
		}(i)
	}
	close(start)
	wg.Wait()
}
