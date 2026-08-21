package domain

import "sort"

type WeightRule struct {
	ID      string
	Target  string
	Percent int
}

func NormalizeWeights(items []WeightRule) []WeightRule {
	out := append([]WeightRule(nil), items...)
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	total := 0
	for _, v := range out {
		if v.Percent > 0 {
			total += v.Percent
		}
	}
	if total == 0 {
		return out
	}
	for i := range out {
		out[i].Percent = out[i].Percent * 100 / total
	}
	return out
}
func PickWeight(items []WeightRule, bucket int) WeightRule {
	sum := 0
	for _, v := range NormalizeWeights(items) {
		sum += v.Percent
		if bucket%100 < sum {
			return v
		}
	}
	return WeightRule{}
}
