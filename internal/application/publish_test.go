package application

import (
	"context"
	"example.com/api-traffic-governance/internal/domain"
	"example.com/api-traffic-governance/internal/xds"
	"sync"
	"testing"
)

func TestPublisherRejectsBlankRoute(t *testing.T) {
	p := NewPublisher(xds.New("secret"))
	if _, err := p.Publish(context.Background(), []domain.Route{{}}, nil); err == nil {
		t.Fatal("blank route was published")
	}
	if _, err := p.Publish(context.Background(), nil, []domain.Quota{{}}); err == nil {
		t.Fatal("blank quota was published")
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := p.Publish(ctx, []domain.Route{{ID: "r", ServiceID: "s", UpstreamIDs: []string{"u"}}}, nil); err == nil {
		t.Fatal("canceled publish succeeded")
	}
	if p.Revision() != 0 {
		t.Fatalf("revision advanced after rejection: %d", p.Revision())
	}
}

func TestPublisherRevisionConcurrent(t *testing.T) {
	p := NewPublisher(xds.New("secret"))
	start := make(chan struct{})
	var wg sync.WaitGroup
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			_, _ = p.Publish(context.Background(), []domain.Route{{ID: "r", ServiceID: "s", UpstreamIDs: []string{"u"}}}, nil)
		}()
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	wg.Add(1)
	go func() {
		defer wg.Done()
		_, _ = p.Publish(ctx, []domain.Route{{ID: "r-cancel", ServiceID: "s", UpstreamIDs: []string{"u"}}}, nil)
	}()
	close(start)
	wg.Wait()
	if p.Revision() != 3 {
		t.Fatalf("revision = %d", p.Revision())
	}
}

func TestPublisherRejectsRevisionOverflow(t *testing.T) {
	p := NewPublisher(xds.New("secret"))
	p.revision = ^uint64(0)
	if _, err := p.Publish(context.Background(), []domain.Route{{ID: "r", ServiceID: "s", UpstreamIDs: []string{"u"}}}, nil); err == nil {
		t.Fatal("revision overflow was accepted")
	}
}
