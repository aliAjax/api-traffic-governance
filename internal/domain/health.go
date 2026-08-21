package domain

import (
	"fmt"
	"time"
)

type HealthEvent struct {
	UpstreamID string
	Healthy    bool
	LatencyMs  int
	StatusCode int
	At         time.Time
}

func (e HealthEvent) Validate() error {
	if e.UpstreamID == "" || e.LatencyMs < 0 || e.StatusCode < 0 {
		return fmt.Errorf("%w: health event", ErrInvalidInput)
	}
	if e.At.IsZero() {
		return fmt.Errorf("%w: event time", ErrInvalidInput)
	}
	return nil
}

type HealthWindow struct {
	Events []HealthEvent
	Max    int
}

func (w *HealthWindow) Add(e HealthEvent) {
	if w.Max < 1 {
		w.Max = 100
	}
	if len(w.Events) >= w.Max {
		w.Events = w.Events[1:]
	}
	w.Events = append(w.Events, e)
}
func (w HealthWindow) FailureRate() float64 {
	if len(w.Events) == 0 {
		return 0
	}
	bad := 0
	for _, e := range w.Events {
		if !e.Healthy || e.StatusCode >= 500 {
			bad++
		}
	}
	return float64(bad) / float64(len(w.Events))
}
