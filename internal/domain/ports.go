package domain

import "context"

type ServiceRepository interface {
	CreateService(context.Context, Service) error
	GetService(context.Context, string) (Service, error)
	ListServices(context.Context, string) ([]Service, error)
}
type UpstreamRepository interface {
	CreateUpstream(context.Context, Upstream) error
	GetUpstream(context.Context, string) (Upstream, error)
	ListUpstreams(context.Context, string) ([]Upstream, error)
	UpdateHealth(context.Context, string, bool, int) error
}
type RouteRepository interface {
	CreateRoute(context.Context, Route) error
	GetRoute(context.Context, string) (Route, error)
	ListRoutes(context.Context, string) ([]Route, error)
}
type QuotaRepository interface {
	CreateQuota(context.Context, Quota) error
	GetQuota(context.Context, string) (Quota, error)
	ListQuotas(context.Context, string) ([]Quota, error)
	UpdateQuota(context.Context, Quota) error
}
type ExperimentRepository interface {
	CreateExperiment(context.Context, Experiment) error
	GetExperiment(context.Context, string) (Experiment, error)
	ListExperiments(context.Context, string) ([]Experiment, error)
	UpdateExperiment(context.Context, Experiment) error
}
type EventSink interface {
	Publish(context.Context, string, any) error
}
