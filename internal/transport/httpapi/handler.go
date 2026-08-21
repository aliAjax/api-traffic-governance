package httpapi

import (
	"encoding/json"
	"errors"
	"example.com/api-traffic-governance/internal/application"
	"example.com/api-traffic-governance/internal/domain"
	"example.com/api-traffic-governance/internal/fault"
	"example.com/api-traffic-governance/internal/limiter"
	"example.com/api-traffic-governance/internal/repository"
	"example.com/api-traffic-governance/internal/routing"
	"example.com/api-traffic-governance/internal/xds"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Handler struct {
	svc      *application.Service
	store    *repository.Store
	faults   *fault.Manager
	selector *routing.Selector
	limits   map[string]*limiter.Bucket
	xds      *xds.Manager
}

func New(s *application.Service, st *repository.Store, f *fault.Manager, x *xds.Manager) *Handler {
	return &Handler{svc: s, store: st, faults: f, selector: routing.New(), limits: map[string]*limiter.Bucket{}, xds: x}
}
func (h *Handler) Routes() *http.ServeMux {
	m := http.NewServeMux()
	m.HandleFunc("/healthz", h.health)
	m.HandleFunc("/readyz", h.health)
	m.HandleFunc("/api/v1/services", h.services)
	m.HandleFunc("/api/v1/upstreams", h.upstreams)
	m.HandleFunc("/api/v1/routes", h.routes)
	m.HandleFunc("/api/v1/quotas", h.quotas)
	m.HandleFunc("/api/v1/experiments", h.experiments)
	m.HandleFunc("/api/v1/experiments/", h.experimentDetail)
	m.HandleFunc("/api/v1/decision/explain", h.explain)
	m.HandleFunc("/api/v1/snapshots", h.snapshots)
	return m
}
func (h *Handler) health(w http.ResponseWriter, _ *http.Request) {
	write(w, 200, map[string]any{"status": "ok", "service": "api-traffic-governance"})
}
func (h *Handler) services(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		v, e := h.store.ListServices(r.Context(), r.URL.Query().Get("tenant_id"))
		respond(w, v, e)
		return
	}
	if r.Method != "POST" {
		w.WriteHeader(405)
		return
	}
	var v domain.Service
	if !decode(r, &v) {
		writeErr(w, domain.ErrInvalidInput)
		return
	}
	out, e := h.svc.CreateService(r.Context(), v)
	if e != nil {
		writeErr(w, e)
		return
	}
	write(w, 201, out)
}
func (h *Handler) upstreams(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		v, e := h.store.ListUpstreams(r.Context(), r.URL.Query().Get("service_id"))
		respond(w, v, e)
		return
	}
	if r.Method != "POST" {
		w.WriteHeader(405)
		return
	}
	var v domain.Upstream
	if !decode(r, &v) {
		writeErr(w, domain.ErrInvalidInput)
		return
	}
	out, e := h.svc.CreateUpstream(r.Context(), v)
	if e != nil {
		writeErr(w, e)
		return
	}
	write(w, 201, out)
}
func (h *Handler) routes(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		v, e := h.store.ListRoutes(r.Context(), r.URL.Query().Get("service_id"))
		respond(w, v, e)
		return
	}
	if r.Method != "POST" {
		w.WriteHeader(405)
		return
	}
	var v domain.Route
	if !decode(r, &v) {
		writeErr(w, domain.ErrInvalidInput)
		return
	}
	out, e := h.svc.CreateRoute(r.Context(), v)
	if e != nil {
		writeErr(w, e)
		return
	}
	write(w, 201, out)
}
func (h *Handler) quotas(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		v, e := h.store.ListQuotas(r.Context(), r.URL.Query().Get("tenant_id"))
		respond(w, v, e)
		return
	}
	if r.Method != "POST" {
		w.WriteHeader(405)
		return
	}
	var v domain.Quota
	if !decode(r, &v) {
		writeErr(w, domain.ErrInvalidInput)
		return
	}
	out, e := h.svc.CreateQuota(r.Context(), v)
	if e != nil {
		writeErr(w, e)
		return
	}
	h.limits[out.ID] = limiter.New(out.Burst, out.RatePerSecond)
	write(w, 201, out)
}
func (h *Handler) experiments(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		v, e := h.store.ListExperiments(r.Context(), r.URL.Query().Get("tenant_id"))
		respond(w, v, e)
		return
	}
	if r.Method != "POST" {
		w.WriteHeader(405)
		return
	}
	var v domain.Experiment
	if !decode(r, &v) {
		writeErr(w, domain.ErrInvalidInput)
		return
	}
	out, e := h.svc.CreateExperiment(r.Context(), v)
	if e != nil {
		writeErr(w, e)
		return
	}
	_ = h.faults.Put(out)
	write(w, 201, out)
}
func (h *Handler) experimentDetail(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/v1/experiments/")
	if r.Method != "POST" || !strings.HasSuffix(r.URL.Path, "/revoke") {
		w.WriteHeader(405)
		return
	}
	id = strings.TrimSuffix(strings.TrimPrefix(r.URL.Path, "/api/v1/experiments/"), "/revoke")
	if e := h.faults.Revoke(id); e != nil {
		writeErr(w, e)
		return
	}
	write(w, 200, map[string]string{"id": id, "status": "expired"})
}
func (h *Handler) explain(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(405)
		return
	}
	var input struct {
		Path    string            `json:"path"`
		Method  string            `json:"method"`
		Headers map[string]string `json:"headers"`
		Query   map[string]string `json:"query"`
	}
	if !decode(r, &input) {
		writeErr(w, domain.ErrInvalidInput)
		return
	}
	if input.Path == "" {
		input.Path = r.Header.Get("X-Original-Path")
	}
	if input.Path == "" {
		input.Path = "/"
	}
	if input.Method == "" {
		input.Method = http.MethodGet
	}
	query := url.Values{}
	for k, v := range input.Query {
		query.Set(k, v)
	}
	target, err := http.NewRequestWithContext(r.Context(), input.Method, "http://decision.local"+input.Path+"?"+query.Encode(), nil)
	if err != nil {
		writeErr(w, domain.ErrInvalidInput)
		return
	}
	for k, v := range input.Headers {
		target.Header.Set(k, v)
	}
	tenant := r.Header.Get("X-Tenant-ID")
	key := r.Header.Get("X-Rate-Key")
	if key == "" {
		key = tenant
	}
	for id, b := range h.limits {
		if strings.HasPrefix(id, tenant+"/") && !b.Allow(1) {
			write(w, 429, domain.Decision{Allowed: false, Status: 429, Reason: "rate_limit", PolicyVersion: h.xds.Current().Version})
			return
		}
	}
	routes, _ := h.store.ListRoutes(r.Context(), "")
	for _, route := range routes {
		if routing.Match(route, target) {
			ups, _ := h.store.ListUpstreams(r.Context(), route.ServiceID)
			u, ok := h.selector.Select(key, ups)
			if !ok {
				write(w, 503, domain.Decision{Allowed: false, Status: 503, Reason: "no_healthy_upstream", RouteID: route.ID, PolicyVersion: h.xds.Current().Version})
				return
			}
			write(w, 200, domain.Decision{Allowed: true, Status: 200, Reason: "matched", RouteID: route.ID, UpstreamID: u.ID, PolicyVersion: h.xds.Current().Version})
			return
		}
	}
	write(w, 404, domain.Decision{Allowed: false, Status: 404, Reason: "no_route", PolicyVersion: h.xds.Current().Version})
}
func (h *Handler) snapshots(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		write(w, 200, h.xds.Current())
		return
	}
	routes, _ := h.store.ListRoutes(r.Context(), "")
	quotas, _ := h.store.ListQuotas(r.Context(), "")
	write(w, 201, h.xds.Compile(routes, quotas))
}
func decode(r *http.Request, v any) bool {
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(v) == nil
}
func write(w http.ResponseWriter, s int, v any) {
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(s)
	_ = json.NewEncoder(w).Encode(v)
}
func respond(w http.ResponseWriter, v any, e error) {
	if e != nil {
		writeErr(w, e)
		return
	}
	write(w, 200, v)
}
func writeErr(w http.ResponseWriter, e error) {
	s := 500
	if e == domain.ErrInvalidInput {
		s = 400
	}
	if e == domain.ErrNotFound {
		s = 404
	}
	if e == domain.ErrConflict {
		s = 409
	}
	if errors.Is(e, domain.ErrUnavailable) {
		s = 503
	}
	write(w, s, map[string]string{"error": e.Error()})
}

var _ = contextTimeout

func contextTimeout() time.Duration { return time.Second }
