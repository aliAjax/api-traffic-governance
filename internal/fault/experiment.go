package fault

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"example.com/api-traffic-governance/internal/domain"
)

type Manager struct {
	mu    sync.RWMutex
	items map[string]domain.Experiment
}

func New() *Manager { return &Manager{items: map[string]domain.Experiment{}} }
func (m *Manager) Put(e domain.Experiment) error {
	if e.ExpiresAt.IsZero() {
		e.ExpiresAt = time.Now().Add(time.Hour)
	}
	if err := e.Validate(); err != nil {
		return err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.items[e.ID] = e
	return nil
}
func (m *Manager) Revoke(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	e, ok := m.items[id]
	if !ok {
		return domain.ErrNotFound
	}
	e.Status = domain.ExperimentExpired
	m.items[id] = e
	return nil
}
func (m *Manager) Apply(ctx context.Context, r *http.Request, tenant string) (int, string, error) {
	m.mu.RLock()
	var matched *domain.Experiment
	for _, e := range m.items {
		if e.TenantID == tenant && e.Status == domain.ExperimentActive && e.ExpiresAt.After(time.Now()) && e.Scope == r.URL.Path {
			ee := e
			matched = &ee
			break
		}
	}
	m.mu.RUnlock()
	if matched == nil {
		return 0, "", nil
	}
	if ctx.Err() != nil {
		return 0, "", ctx.Err()
	}
	if matched.DelayMs > 0 {
		select {
		case <-time.After(time.Duration(matched.DelayMs) * time.Millisecond):
		case <-ctx.Done():
			return 0, "injected delay", nil
		}
		return 0, "injected delay", nil
	}
	return 0, "", nil
}
func (m *Manager) List(tenant string) []domain.Experiment {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := []domain.Experiment{}
	for _, e := range m.items {
		if tenant == "" || e.TenantID == tenant {
			out = append(out, e)
		}
	}
	return out
}
func ValidateScope(scope string) error {
	if scope == "" || scope[0] != '/' {
		return fmt.Errorf("%w: scope", domain.ErrInvalidInput)
	}
	return nil
}
