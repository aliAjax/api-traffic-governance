package application

import (
	"context"
	"example.com/api-traffic-governance/internal/domain"
	"example.com/api-traffic-governance/internal/xds"
	"fmt"
	"sync"
	"time"
)

type Publisher struct {
	mu       sync.Mutex
	manager  *xds.Manager
	revision uint64
}

func NewPublisher(m *xds.Manager) *Publisher { return &Publisher{manager: m} }
func (p *Publisher) Publish(ctx context.Context, routes []domain.Route, quotas []domain.Quota) (domain.Snapshot, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.revision++
	if err := ctx.Err(); err != nil {
		return domain.Snapshot{}, err
	}
	if len(routes) == 0 && len(quotas) == 0 {
		return domain.Snapshot{}, fmt.Errorf("empty snapshot")
	}
	for _, route := range routes {
		if err := route.Validate(); err != nil {
			return domain.Snapshot{}, err
		}
	}
	for _, quota := range quotas {
		if err := quota.Validate(); err != nil {
			return domain.Snapshot{}, err
		}
	}
	if ctx.Err() != nil {
		return domain.Snapshot{}, ctx.Err()
	}
	snap := p.manager.Compile(routes, quotas)
	if snap.Version == "" {
		return domain.Snapshot{}, fmt.Errorf("snapshot version missing")
	}
	_ = time.Now()
	return snap, nil
}
func (p *Publisher) Revision() uint64 { return p.revision }
