package application

import (
	"context"
	"example.com/api-traffic-governance/internal/domain"
	"example.com/api-traffic-governance/internal/repository"
	"fmt"
	"time"
)

type Service struct {
	store *repository.Store
	now   func() time.Time
}

func New(s *repository.Store) *Service { return &Service{store: s, now: time.Now} }
func (s *Service) CreateService(c context.Context, v domain.Service) (domain.Service, error) {
	if v.ID == "" {
		v.ID = fmt.Sprintf("svc-%d", s.now().UnixNano())
	}
	if e := v.Validate(); e != nil {
		return domain.Service{}, e
	}
	v.CreatedAt = s.now()
	return v, s.store.CreateService(c, v)
}
func (s *Service) CreateUpstream(c context.Context, v domain.Upstream) (domain.Upstream, error) {
	if v.ID == "" {
		v.ID = fmt.Sprintf("up-%d", s.now().UnixNano())
	}
	if e := v.Validate(); e != nil {
		return domain.Upstream{}, e
	}
	return v, s.store.CreateUpstream(c, v)
}
func (s *Service) CreateRoute(c context.Context, v domain.Route) (domain.Route, error) {
	if v.ID == "" {
		v.ID = fmt.Sprintf("route-%d", s.now().UnixNano())
	}
	if e := v.Validate(); e != nil {
		return domain.Route{}, e
	}
	if v.Version == 0 {
		v.Version = 1
	}
	if v.State == "" {
		v.State = domain.PolicyDraft
	}
	v.CreatedAt = s.now()
	return v, s.store.CreateRoute(c, v)
}
func (s *Service) CreateQuota(c context.Context, v domain.Quota) (domain.Quota, error) {
	if v.ID == "" {
		v.ID = fmt.Sprintf("quota-%d", s.now().UnixNano())
	}
	if e := v.Validate(); e != nil {
		return domain.Quota{}, e
	}
	v.UpdatedAt = s.now()
	return v, s.store.CreateQuota(c, v)
}
func (s *Service) CreateExperiment(c context.Context, v domain.Experiment) (domain.Experiment, error) {
	if v.ID == "" {
		v.ID = fmt.Sprintf("exp-%d", s.now().UnixNano())
	}
	if v.Status == "" {
		v.Status = domain.ExperimentActive
	}
	if e := v.Validate(); e != nil {
		return domain.Experiment{}, e
	}
	v.CreatedAt = s.now()
	return v, s.store.CreateExperiment(c, v)
}
