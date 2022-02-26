package grpcHealth

import (
	"context"
	"go.uber.org/zap"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func HealthCheck(client grpc_health_v1.HealthClient, serviceName string, log *zap.SugaredLogger) {
	req := &grpc_health_v1.HealthCheckRequest{
		Service: serviceName,
	}
	hcResponse, err := client.Check(context.Background(), req)
	if err != nil || hcResponse.Status != grpc_health_v1.HealthCheckResponse_SERVING {
		log.Errorw("Health Check Error", "error", err, "response", hcResponse)
	} else {
		log.Debugw("Health Check passing", "response", hcResponse)
	}
}
