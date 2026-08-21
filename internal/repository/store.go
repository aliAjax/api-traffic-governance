package repository

import (
	"context"
	"encoding/json"
	"example.com/api-traffic-governance/internal/domain"
	"os"
	"path/filepath"
	"sync"
)

type State struct {
	Services    map[string]domain.Service    `json:"services"`
	Upstreams   map[string]domain.Upstream   `json:"upstreams"`
	Routes      map[string]domain.Route      `json:"routes"`
	Quotas      map[string]domain.Quota      `json:"quotas"`
	Experiments map[string]domain.Experiment `json:"experiments"`
}
type Store struct {
	mu   sync.RWMutex
	path string
	s    State
}

func New(path string) (*Store, error) {
	st := &Store{path: path, s: State{map[string]domain.Service{}, map[string]domain.Upstream{}, map[string]domain.Route{}, map[string]domain.Quota{}, map[string]domain.Experiment{}}}
	if b, e := os.ReadFile(path); e == nil {
		if e = json.Unmarshal(b, &st.s); e != nil {
			return nil, e
		}
	} else if !os.IsNotExist(e) {
		return nil, e
	}
	return st, nil
}
func (st *Store) save() error {
	b, e := json.MarshalIndent(st.s, "", "  ")
	if e != nil {
		return e
	}
	if e = os.MkdirAll(filepath.Dir(st.path), 0755); e != nil {
		return e
	}
	tmp := st.path + ".tmp"
	if e = os.WriteFile(tmp, b, 0644); e != nil {
		return e
	}
	return os.Rename(tmp, st.path)
}
func (st *Store) CreateService(_ context.Context, v domain.Service) error {
	st.mu.Lock()
	defer st.mu.Unlock()
	if _, ok := st.s.Services[v.ID]; ok {
		return domain.ErrConflict
	}
	st.s.Services[v.ID] = v
	return st.save()
}
func (st *Store) GetService(_ context.Context, id string) (domain.Service, error) {
	st.mu.RLock()
	defer st.mu.RUnlock()
	v, ok := st.s.Services[id]
	if !ok {
		return domain.Service{}, domain.ErrNotFound
	}
	return v, nil
}
func (st *Store) ListServices(_ context.Context, t string) ([]domain.Service, error) {
	st.mu.RLock()
	defer st.mu.RUnlock()
	o := []domain.Service{}
	for _, v := range st.s.Services {
		if t == "" || v.TenantID == t {
			o = append(o, v)
		}
	}
	return o, nil
}
func (st *Store) CreateUpstream(_ context.Context, v domain.Upstream) error {
	st.mu.Lock()
	defer st.mu.Unlock()
	if _, ok := st.s.Upstreams[v.ID]; ok {
		return domain.ErrConflict
	}
	st.s.Upstreams[v.ID] = v
	return st.save()
}
func (st *Store) GetUpstream(_ context.Context, id string) (domain.Upstream, error) {
	st.mu.RLock()
	defer st.mu.RUnlock()
	v, ok := st.s.Upstreams[id]
	if !ok {
		return domain.Upstream{}, domain.ErrNotFound
	}
	return v, nil
}
func (st *Store) ListUpstreams(_ context.Context, sid string) ([]domain.Upstream, error) {
	st.mu.RLock()
	defer st.mu.RUnlock()
	o := []domain.Upstream{}
	for _, v := range st.s.Upstreams {
		if sid == "" || v.ServiceID == sid {
			o = append(o, v)
		}
	}
	return o, nil
}
func (st *Store) UpdateHealth(_ context.Context, id string, h bool, lat int) error {
	st.mu.Lock()
	defer st.mu.Unlock()
	v, ok := st.s.Upstreams[id]
	if !ok {
		return domain.ErrNotFound
	}
	v.Healthy = h
	v.LatencyMs = lat
	st.s.Upstreams[id] = v
	return st.save()
}
func (st *Store) CreateRoute(_ context.Context, v domain.Route) error {
	st.mu.Lock()
	defer st.mu.Unlock()
	if _, ok := st.s.Routes[v.ID]; ok {
		return domain.ErrConflict
	}
	st.s.Routes[v.ID] = v
	return st.save()
}
func (st *Store) GetRoute(_ context.Context, id string) (domain.Route, error) {
	st.mu.RLock()
	defer st.mu.RUnlock()
	v, ok := st.s.Routes[id]
	if !ok {
		return domain.Route{}, domain.ErrNotFound
	}
	return v, nil
}
func (st *Store) ListRoutes(_ context.Context, sid string) ([]domain.Route, error) {
	st.mu.RLock()
	defer st.mu.RUnlock()
	o := []domain.Route{}
	for _, v := range st.s.Routes {
		if sid == "" || v.ServiceID == sid {
			o = append(o, v)
		}
	}
	return o, nil
}
func (st *Store) CreateQuota(_ context.Context, v domain.Quota) error {
	st.mu.Lock()
	defer st.mu.Unlock()
	if _, ok := st.s.Quotas[v.ID]; ok {
		return domain.ErrConflict
	}
	st.s.Quotas[v.ID] = v
	return st.save()
}
func (st *Store) GetQuota(_ context.Context, id string) (domain.Quota, error) {
	st.mu.RLock()
	defer st.mu.RUnlock()
	v, ok := st.s.Quotas[id]
	if !ok {
		return domain.Quota{}, domain.ErrNotFound
	}
	return v, nil
}
func (st *Store) ListQuotas(_ context.Context, t string) ([]domain.Quota, error) {
	st.mu.RLock()
	defer st.mu.RUnlock()
	o := []domain.Quota{}
	for _, v := range st.s.Quotas {
		if t == "" || v.TenantID == t {
			o = append(o, v)
		}
	}
	return o, nil
}
func (st *Store) UpdateQuota(_ context.Context, v domain.Quota) error {
	st.mu.Lock()
	defer st.mu.Unlock()
	if _, ok := st.s.Quotas[v.ID]; !ok {
		return domain.ErrNotFound
	}
	st.s.Quotas[v.ID] = v
	return st.save()
}
func (st *Store) CreateExperiment(_ context.Context, v domain.Experiment) error {
	st.mu.Lock()
	defer st.mu.Unlock()
	if _, ok := st.s.Experiments[v.ID]; ok {
		return domain.ErrConflict
	}
	st.s.Experiments[v.ID] = v
	return st.save()
}
func (st *Store) GetExperiment(_ context.Context, id string) (domain.Experiment, error) {
	st.mu.RLock()
	defer st.mu.RUnlock()
	v, ok := st.s.Experiments[id]
	if !ok {
		return domain.Experiment{}, domain.ErrNotFound
	}
	return v, nil
}
func (st *Store) ListExperiments(_ context.Context, t string) ([]domain.Experiment, error) {
	st.mu.RLock()
	defer st.mu.RUnlock()
	o := []domain.Experiment{}
	for _, v := range st.s.Experiments {
		if t == "" || v.TenantID == t {
			o = append(o, v)
		}
	}
	return o, nil
}
func (st *Store) UpdateExperiment(_ context.Context, v domain.Experiment) error {
	st.mu.Lock()
	defer st.mu.Unlock()
	if _, ok := st.s.Experiments[v.ID]; !ok {
		return domain.ErrNotFound
	}
	st.s.Experiments[v.ID] = v
	return st.save()
}
