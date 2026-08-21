package observability

import "sync/atomic"

type Metrics struct {
	Requests atomic.Uint64
	Allowed  atomic.Uint64
	Limited  atomic.Uint64
	Failures atomic.Uint64
}

func (m *Metrics) Snapshot() map[string]uint64 {
	return map[string]uint64{"requests": m.Requests.Load(), "allowed": m.Allowed.Load(), "limited": m.Limited.Load(), "failures": m.Failures.Load()}
}
