package mirror

import (
	"sort"
	"time"
)

func Expired(items []Message, now time.Time) []Message {
	out := items[:0]
	for _, v := range items {
		if !v.ExpiresAt.IsZero() && v.ExpiresAt.Before(now) {
			out = append(out, v)
		}
	}
	return out
}
func SortNewest(items []Message) {
	sort.SliceStable(items, func(i, j int) bool { return items[i].CreatedAt.After(items[j].CreatedAt) })
}
