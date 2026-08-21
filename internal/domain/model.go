package domain

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

var (
	ErrInvalidInput = errors.New("invalid input")
	ErrNotFound     = errors.New("not found")
	ErrConflict     = errors.New("conflict")
	ErrUnavailable  = errors.New("unavailable")
)

type State string

const (
	PolicyDraft       State = "draft"
	PolicyPublished   State = "published"
	PolicyArchived    State = "archived"
	ExperimentActive  State = "active"
	ExperimentExpired State = "expired"
	CircuitClosed     State = "closed"
	CircuitOpen       State = "open"
	CircuitHalfOpen   State = "half_open"
)

type Service struct {
	ID        string    `json:"id"`
	TenantID  string    `json:"tenant_id"`
	Name      string    `json:"name"`
	Region    string    `json:"region"`
	CreatedAt time.Time `json:"created_at"`
}
type Upstream struct {
	ID        string `json:"id"`
	ServiceID string `json:"service_id"`
	URL       string `json:"url"`
	Weight    int    `json:"weight"`
	Region    string `json:"region"`
	Healthy   bool   `json:"healthy"`
	LatencyMs int    `json:"latency_ms"`
}
type Match struct {
	PathPrefix string            `json:"path_prefix"`
	Methods    []string          `json:"methods"`
	Headers    map[string]string `json:"headers"`
	Query      map[string]string `json:"query"`
}
type Route struct {
	ID          string    `json:"id"`
	ServiceID   string    `json:"service_id"`
	Match       Match     `json:"match"`
	UpstreamIDs []string  `json:"upstream_ids"`
	Version     int       `json:"version"`
	State       State     `json:"state"`
	TimeoutMs   int       `json:"timeout_ms"`
	Retries     int       `json:"retries"`
	CreatedAt   time.Time `json:"created_at"`
}
type Quota struct {
	ID            string    `json:"id"`
	TenantID      string    `json:"tenant_id"`
	Key           string    `json:"key"`
	RatePerSecond int       `json:"rate_per_second"`
	Burst         int       `json:"burst"`
	Concurrent    int       `json:"concurrent"`
	Used          int       `json:"used"`
	UpdatedAt     time.Time `json:"updated_at"`
}
type Experiment struct {
	ID        string    `json:"id"`
	TenantID  string    `json:"tenant_id"`
	Name      string    `json:"name"`
	Scope     string    `json:"scope"`
	DelayMs   int       `json:"delay_ms"`
	Status    State     `json:"status"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}
type Snapshot struct {
	Version   string    `json:"version"`
	Nonce     string    `json:"nonce"`
	Routes    []Route   `json:"routes"`
	Quotas    []Quota   `json:"quotas"`
	CreatedAt time.Time `json:"created_at"`
}
type Decision struct {
	Allowed       bool   `json:"allowed"`
	Status        int    `json:"status"`
	Reason        string `json:"reason"`
	RouteID       string `json:"route_id,omitempty"`
	UpstreamID    string `json:"upstream_id,omitempty"`
	PolicyVersion string `json:"policy_version"`
	LatencyMs     int64  `json:"latency_ms"`
}

func (s Service) Validate() error {
	if strings.TrimSpace(s.ID) == "" || strings.TrimSpace(s.TenantID) == "" {
		return fmt.Errorf("%w: service fields", ErrInvalidInput)
	}
	return nil
}
func (u Upstream) Validate() error {
	if u.ID == "" || u.ServiceID == "" || u.Weight < 0 {
		return fmt.Errorf("%w: upstream fields", ErrInvalidInput)
	}
	return nil
}
func (r Route) Validate() error {
	if r.ID == "" || r.ServiceID == "" || len(r.UpstreamIDs) == 0 {
		return fmt.Errorf("%w: route fields", ErrInvalidInput)
	}
	if r.TimeoutMs < 0 || r.Retries < 0 || r.Retries > 10 {
		return fmt.Errorf("%w: route limits", ErrInvalidInput)
	}
	return nil
}
func (q Quota) Validate() error {
	if q.TenantID == "" || q.Key == "" || q.RatePerSecond <= 0 || q.Burst <= 0 {
		return fmt.Errorf("%w: quota limits", ErrInvalidInput)
	}
	return nil
}
func (e Experiment) Validate() error {
	if e.TenantID == "" || e.Name == "" || e.Scope == "" || e.DelayMs < 0 {
		return fmt.Errorf("%w: experiment fields", ErrInvalidInput)
	}
	if e.ExpiresAt.Before(time.Now()) {
		return fmt.Errorf("%w: experiment expired", ErrInvalidInput)
	}
	return nil
}
