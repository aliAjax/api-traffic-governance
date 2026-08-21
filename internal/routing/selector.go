package routing

import (
	"example.com/api-traffic-governance/internal/domain"
	"hash/fnv"
	"sort"
	"sync"
)

type Selector struct {
	mu     sync.Mutex
	cursor map[string]int
}

func New() *Selector { return &Selector{cursor: map[string]int{}} }
func (s *Selector) Select(key string, items []domain.Upstream) (domain.Upstream, bool) {
	eligible := make([]domain.Upstream, 0, len(items))
	for _, v := range items {
		if v.Healthy && v.Weight > 0 {
			eligible = append(eligible, v)
		}
	}
	if len(eligible) == 0 {
		return domain.Upstream{}, false
	}
	sort.Slice(eligible, func(i, j int) bool { return eligible[i].ID < eligible[j].ID })
	if key != "" {
		h := fnv.New32a()
		_, _ = h.Write([]byte(key))
		return eligible[int(h.Sum32())%len(eligible)], true
	}
	s.mu.Lock()
	i := s.cursor[eligible[0].ServiceID] % len(eligible)
	s.cursor[eligible[0].ServiceID]++
	s.mu.Unlock()
	return eligible[i], true
}
func Weighted(items []domain.Upstream) []domain.Upstream {
	out := []domain.Upstream{}
	for _, v := range items {
		n := v.Weight
		if n > 100 {
			n = 100
		}
		for i := 0; i < n; i++ {
			out = append(out, v)
		}
	}
	return out
}
