package xds

import (
	"example.com/api-traffic-governance/internal/domain"
	"reflect"
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
	if !reflect.DeepEqual(old.Quotas, next.Quotas) {
		if len(next.Quotas) > len(old.Quotas) {
			d.QuotasAdded = len(next.Quotas) - len(old.Quotas)
		} else {
			d.QuotasRemoved = len(old.Quotas) - len(next.Quotas)
		}
	}
	return d
}
