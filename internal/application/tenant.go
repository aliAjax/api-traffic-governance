package application

import (
	"fmt"
	"sync"
	"time"
)

type TenantRegistry struct {
	mu    sync.RWMutex
	items map[string]time.Time
}

func NewTenantRegistry() *TenantRegistry { return &TenantRegistry{items: map[string]time.Time{}} }
func (t *TenantRegistry) Register(id string) error {
	if id == "" {
		return fmt.Errorf("tenant required")
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	if _, ok := t.items[id]; ok {
		return fmt.Errorf("tenant exists")
	}
	t.items[id] = time.Now()
	return nil
}
func (t *TenantRegistry) Exists(id string) bool {
	t.mu.RLock()
	defer t.mu.RUnlock()
	_, ok := t.items[id]
	return ok
}
func (t *TenantRegistry) List() []string {
	t.mu.RLock()
	defer t.mu.RUnlock()
	out := []string{}
	for id := range t.items {
		out = append(out, id)
	}
	return out
}
