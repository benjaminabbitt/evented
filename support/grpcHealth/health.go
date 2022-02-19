package grpcHealth

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func RegisterHealthChecks(rpc *grpc.Server, name string) *HealthReporter {
	hlth := health.NewServer()
	grpc_health_v1.RegisterHealthServer(rpc, hlth)
	return &HealthReporter{
		name:       name,
		hlthServer: hlth,
	}
}

type HealthReporter struct {
	name       string
	hlthServer *health.Server
}

func (o *HealthReporter) SetStatus(status grpc_health_v1.HealthCheckResponse_ServingStatus) {
	o.hlthServer.SetServingStatus(o.name, status)
}

func (o *HealthReporter) Shutdown() {
	o.hlthServer.Shutdown()
}

func (o *HealthReporter) OK() {
	o.hlthServer.SetServingStatus(o.name, grpc_health_v1.HealthCheckResponse_SERVING)
}
