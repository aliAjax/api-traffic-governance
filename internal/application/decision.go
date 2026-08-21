package application

import (
	"example.com/api-traffic-governance/internal/domain"
	"example.com/api-traffic-governance/internal/routing"
	"net/http"
)

type DecisionEngine struct {
	routes    []domain.Route
	upstreams map[string][]domain.Upstream
	selector  *routing.Selector
}

func NewDecisionEngine(routes []domain.Route, ups map[string][]domain.Upstream) *DecisionEngine {
	return &DecisionEngine{routes: routes, upstreams: ups, selector: routing.New()}
}
func (e *DecisionEngine) Evaluate(req *http.Request, key string) domain.Decision {
	for _, r := range e.routes {
		if !routing.Match(r, req) {
			continue
		}
		u, ok := e.selector.Select(key, e.upstreams[r.ServiceID])
		if !ok {
			return domain.Decision{Allowed: false, Status: 503, Reason: "no_healthy_upstream", RouteID: r.ID}
		}
		return domain.Decision{Allowed: true, Status: 200, Reason: "matched", RouteID: r.ID, UpstreamID: u.ID}
	}
	return domain.Decision{Allowed: false, Status: 404, Reason: "no_route"}
}
