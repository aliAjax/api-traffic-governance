package routing

import (
	"example.com/api-traffic-governance/internal/domain"
	"sort"
)

func SortByLatency(items []domain.Upstream) []domain.Upstream {
	out := append([]domain.Upstream(nil), items...)
	sort.SliceStable(out, func(i, j int) bool {
		if out[i].Healthy != out[j].Healthy {
			return out[i].Healthy
		}
		return out[i].LatencyMs < out[j].LatencyMs
	})
	return out
}
func Healthy(items []domain.Upstream) []domain.Upstream {
	out := []domain.Upstream{}
	for _, v := range items {
		if v.Healthy {
			out = append(out, v)
		}
	}
	return out
}
