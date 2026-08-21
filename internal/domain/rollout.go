package domain

import (
	"fmt"
	"sort"
	"time"
)

type Rollout struct {
	ID          string
	TenantID    string
	RouteID     string
	FromVersion int
	ToVersion   int
	Percent     int
	State       State
	StartedAt   time.Time
	FinishedAt  time.Time
}

func (r Rollout) Validate() error {
	if r.ID == "" || r.TenantID == "" || r.RouteID == "" || r.FromVersion < 0 || r.ToVersion < 1 {
		return fmt.Errorf("%w: rollout fields", ErrInvalidInput)
	}
	if r.Percent < 0 || r.Percent > 100 {
		return fmt.Errorf("%w: rollout percent", ErrInvalidInput)
	}
	return nil
}
func SortRollouts(v []Rollout) {
	sort.Slice(v, func(i, j int) bool { return v[i].StartedAt.Before(v[j].StartedAt) })
}
func (r Rollout) Complete(now time.Time) Rollout {
	r.State = PolicyPublished
	r.FinishedAt = now
	return r
}
