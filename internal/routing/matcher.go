package routing

import (
	"example.com/api-traffic-governance/internal/domain"
	"net/http"
	"strings"
)

func Match(r domain.Route, req *http.Request) bool {
	m := r.Match
	if m.PathPrefix != "" && !strings.HasPrefix(req.URL.Path, m.PathPrefix) {
		return false
	}
	if len(m.Methods) > 0 {
		ok := false
		for _, v := range m.Methods {
			if strings.EqualFold(v, req.Method) {
				ok = true
			}
		}
		if !ok {
			return false
		}
	}
	for k, v := range m.Headers {
		if req.Header.Get(k) != v {
			return false
		}
	}
	for k, v := range m.Query {
		if req.URL.Query().Get(k) != v {
			return false
		}
	}
	return true
}
func Conflicts(a, b domain.Route) bool {
	return a.ServiceID == b.ServiceID && a.Match.PathPrefix == b.Match.PathPrefix
}
