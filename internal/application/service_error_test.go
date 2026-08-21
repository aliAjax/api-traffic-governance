package application

import (
	"context"
	"errors"
	"example.com/api-traffic-governance/internal/domain"
	"testing"
)

func TestCreateServicePreservesValidationSentinel(t *testing.T) {
	s := &Service{store: nil}
	_, err := s.CreateService(context.Background(), domain.Service{ID: "svc", TenantID: "tenant"})
	if !errors.Is(err, domain.ErrInvalidInput) {
		t.Fatalf("error = %v", err)
	}
}

func TestCreateUpstreamPreservesValidationSentinel(t *testing.T) {
	s := &Service{store: nil}
	_, err := s.CreateUpstream(context.Background(), domain.Upstream{ID: "u", ServiceID: ""})
	if !errors.Is(err, domain.ErrInvalidInput) {
		t.Fatalf("error = %v", err)
	}
}

func TestCreateRoutePreservesValidationSentinel(t *testing.T) {
	s := &Service{store: nil}
	_, err := s.CreateRoute(context.Background(), domain.Route{ID: "r", ServiceID: "s"})
	if !errors.Is(err, domain.ErrInvalidInput) {
		t.Fatalf("error = %v", err)
	}
}

func TestCreateQuotaPreservesValidationSentinel(t *testing.T) {
	s := &Service{store: nil}
	_, err := s.CreateQuota(context.Background(), domain.Quota{ID: "q", TenantID: "t"})
	if !errors.Is(err, domain.ErrInvalidInput) {
		t.Fatalf("error = %v", err)
	}
}

func TestCreateExperimentPreservesValidationSentinel(t *testing.T) {
	s := &Service{store: nil}
	_, err := s.CreateExperiment(context.Background(), domain.Experiment{ID: "e", TenantID: "t"})
	if !errors.Is(err, domain.ErrInvalidInput) {
		t.Fatalf("error = %v", err)
	}
}
