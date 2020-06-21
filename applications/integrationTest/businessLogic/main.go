package main

import (
	"github.com/benjaminabbitt/evented/applications/integrationTest/businessLogic/businessLogic"
	"github.com/benjaminabbitt/evented/applications/integrationTest/businessLogic/configuration"
	"github.com/benjaminabbitt/evented/proto/evented/business"
	"github.com/benjaminabbitt/evented/support"
	"github.com/benjaminabbitt/evented/support/consul"
	"github.com/benjaminabbitt/evented/support/grpcWithInterceptors"
	"github.com/google/uuid"
	_ "github.com/spf13/viper/remote"
	"github.com/uber/jaeger-client-go"
	"go.uber.org/zap"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

var log *zap.SugaredLogger

func main() {
	log = support.Log()
	defer log.Sync()

	config := configuration.Configuration{}
	config.Initialize(log)

	setupConsul(config)

	tracer, closer := jaeger.NewTracer(config.AppName(),
		jaeger.NewConstSampler(true),
		jaeger.NewInMemoryReporter(),
	)

	defer closer.Close()

	rpc := grpcWithInterceptors.GenerateConfiguredServer(log.Desugar(), tracer)

	server := businessLogic.NewPlaceholderBusinessLogicServer(log)
	business.RegisterBusinessLogicServer(rpc, server)

	health := health.NewServer()
	grpc_health_v1.RegisterHealthServer(rpc, health)
	health.Resume()

	port := config.Port()
	log.Infow("Starting Business Server...", "port", port)
	server.Listen(port, tracer)
}

func setupConsul(config configuration.Configuration) {

	consul := consul.EventedConsul{}
	id, err := uuid.NewRandom()
	if err != nil {
		log.Error(err)
	}
	consul.Register(config.AppName(), id.String(), config.Port())

}
