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
	blocked := map[string]bool{}
	for _, f := range fields {
		blocked[strings.ToLower(f)] = true
	}
	out := make(map[string]any, len(payload))
	for k, v := range payload {
		if blocked[strings.ToLower(k)] {
			out[k] = "[REDACTED]"
		} else {
			out[k] = deepCopy(v)
		}
	}
	return out
}

func deepCopy(v any) any {
	switch val := v.(type) {
	case map[string]any:
		cp := make(map[string]any, len(val))
		for k, vv := range val {
			cp[k] = deepCopy(vv)
		}
		return cp
	case []any:
		cp := make([]any, len(val))
		for i, vv := range val {
			cp[i] = deepCopy(vv)
		}
		return cp
	default:
		return v
	}
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
