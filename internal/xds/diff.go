package xds

import (
	"example.com/api-traffic-governance/internal/domain"
)

type Diff struct {
	RoutesAdded   int
	RoutesRemoved int
	QuotasAdded   int
	QuotasRemoved int
}

func Compare(old, next domain.Snapshot) Diff {
	d := Diff{}
	oldR := map[string]bool{}
	newR := map[string]bool{}
	for _, v := range old.Routes {
		oldR[v.ID] = true
	}
	for _, v := range next.Routes {
		newR[v.ID] = true
		if !oldR[v.ID] {
			d.RoutesAdded++
		}
	}
	for k := range oldR {
		if !newR[k] {
			d.RoutesRemoved++
		}
	}
	oldQ := map[string]domain.Quota{}
	newQ := map[string]domain.Quota{}
	for _, q := range old.Quotas {
		oldQ[q.ID] = q
	}
	for _, q := range next.Quotas {
		newQ[q.ID] = q
	}
	for id, q := range newQ {
		if prev, ok := oldQ[id]; !ok {
			d.QuotasAdded++
		} else if q != prev {
			d.QuotasAdded++
			d.QuotasRemoved++
		}
	}
	for id := range oldQ {
		if _, ok := newQ[id]; !ok {
			d.QuotasRemoved++
		}
	}
	return d
}
