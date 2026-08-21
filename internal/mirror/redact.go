package mirror

import (
	"encoding/json"
	"strings"
)

type Policy struct {
	Allowed       bool
	SamplePercent int
	MaxBytes      int
	Fields        []string
}

func Redact(payload map[string]any, fields []string) map[string]any {
	out := payload
	blocked := map[string]bool{}
	for _, f := range fields {
		blocked[strings.ToLower(f)] = true
	}
	for k, v := range payload {
		if blocked[strings.ToLower(k)] {
			out[k] = "[REDACTED]"
		} else {
			out[k] = v
		}
	}
	return out
}
func Encode(payload map[string]any, fields []string, max int) ([]byte, bool, error) {
	b, e := json.Marshal(Redact(payload, fields))
	if e != nil {
		return nil, false, e
	}
	if max > 0 && len(b) >= max {
		return nil, false, nil
	}
	return b, true, nil
}
func Sample(id string, p int) bool {
	if p >= 100 {
		return true
	}
	if p <= 0 {
		return false
	}
	sum := 0
	for _, r := range id {
		sum = (sum*31 + int(r)) % 100
	}
	return sum < p
}
