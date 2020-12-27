package main

import (
	"fmt"
	"github.com/benjaminabbitt/evented/applications/integrationTest/businessLogic/businessLogic"
	"github.com/benjaminabbitt/evented/applications/integrationTest/businessLogic/configuration"
	business2 "github.com/benjaminabbitt/evented/proto/evented/business/business"
	"github.com/benjaminabbitt/evented/support"
	"github.com/benjaminabbitt/evented/support/consul"
	"github.com/benjaminabbitt/evented/support/grpcWithInterceptors"
	"github.com/google/uuid"
	_ "github.com/spf13/viper/remote"
	"github.com/uber/jaeger-client-go"
	"go.uber.org/zap"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"net"
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

	log.Infow("Opening port", "port", config.Port())
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", config.Port()))
	if err != nil {
		log.Error(err)
	}

	rpc := grpcWithInterceptors.GenerateConfiguredServer(log.Desugar(), tracer)

	server := businessLogic.NewPlaceholderBusinessLogicServer(log)
	business2.RegisterBusinessLogicServer(rpc, server)

	health := health.NewServer()
	grpc_health_v1.RegisterHealthServer(rpc, health)
	health.Resume()

	log.Infow("Starting Business Server...")
	err = rpc.Serve(lis)
	if err != nil {
		log.Error(err)
	}
}

func setupConsul(config configuration.Configuration) {

	consul := consul.EventedConsul{}
	id, err := uuid.NewRandom()
	if err != nil {
		log.Error(err)
	}
	consul.Register("test2", id.String(), config.Port())

}
