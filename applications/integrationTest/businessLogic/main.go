package main

import (
	"github.com/benjaminabbitt/evented/applications/integrationTest/businessLogic/businessLogic"
	"github.com/benjaminabbitt/evented/applications/integrationTest/businessLogic/configuration"
	"github.com/benjaminabbitt/evented/proto/gen/github.com/benjaminabbitt/evented/proto/evented"
	"github.com/benjaminabbitt/evented/support"
	"github.com/benjaminabbitt/evented/support/consul"
	"github.com/benjaminabbitt/evented/support/grpcHealth"
	"github.com/benjaminabbitt/evented/support/grpcWithInterceptors"
	"github.com/google/uuid"
	_ "github.com/spf13/viper/remote"
	"github.com/uber/jaeger-client-go"
	"go.uber.org/zap"
)

var (
	log *zap.SugaredLogger
)

func main() {
	log = support.Log()
	defer log.Sync()
	support.LogStartup(log, "Sample Business Logic")

	config := configuration.Configuration{}
	config.Initialize(log)

	setupConsul(config)

	tracer, closer := jaeger.NewTracer(config.AppName(),
		jaeger.NewConstSampler(true),
		jaeger.NewInMemoryReporter(),
	)
	defer func() {
		if err := closer.Close(); err != nil {
			log.Errorw("Error closing Jaeger", err)
		}
	}()

	lis, err := support.OpenPort(config.Port(), log)

	rpc := grpcWithInterceptors.GenerateConfiguredServer(log.Desugar(), tracer)

	server := businessLogic.NewPlaceholderBusinessLogicServer(log)
	evented.RegisterBusinessLogicServer(rpc, server)

	grpcHealth.RegisterHealthChecks(rpc, config.AppName(), log)

	log.Infow("Starting Business Server...")
	err = rpc.Serve(lis)
	log.Infow("Serving...")
	if err != nil {
		log.Error(err)
	}
}

func setupConsul(config configuration.Configuration) {
	c := consul.NewEventedConsul(config.ConsulHost(), config.Port())
	id, err := uuid.NewRandom()
	if err != nil {
		log.Error(err)
	}
	err = c.Register(config.AppName(), id.String())
	if err != nil {
		log.Error(err)
	}

}
