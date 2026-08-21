package application

import (
	"context"
	"example.com/api-traffic-governance/internal/domain"
	"example.com/api-traffic-governance/internal/xds"
	"testing"
)

func TestPublisherRejectsBlankRoute(t *testing.T) {
	p := NewPublisher(xds.New("secret"))
	if _, err := p.Publish(context.Background(), []domain.Route{{}}, nil); err == nil {
		t.Fatal("blank route was published")
	}
}
