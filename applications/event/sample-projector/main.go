package main

import (
	"github.com/benjaminabbitt/evented/applications/event/sample-projector/configuration"
	"github.com/benjaminabbitt/evented/applications/event/sample-projector/projector"
	"github.com/benjaminabbitt/evented/proto/gen/github.com/benjaminabbitt/evented/proto/evented"
	"github.com/benjaminabbitt/evented/support"
	"github.com/benjaminabbitt/evented/support/grpcHealth"
	"github.com/benjaminabbitt/evented/support/grpcWithInterceptors"
	"github.com/benjaminabbitt/evented/support/jaeger"
	"go.uber.org/zap"
)

var log *zap.SugaredLogger

/*
Placeholder business logic -- used for Projector integration tests
*/
func main() {
	log = support.Log()
	support.LogStartup(log, "Sample Projector")
	defer func() {
		if err := log.Sync(); err != nil {
			log.Errorw("Error syncing logs", err)
		}
	}()

	initConfig := support.ConfigInit{}
	config := &configuration.Configuration{}
	config = initConfig.Initialize(log, config).(*configuration.Configuration)

	tracer, closer := jaeger.SetupJaeger(config.Name, log)
	defer jaeger.CloseJaeger(closer, log)

	rpc := grpcWithInterceptors.GenerateConfiguredServer(log.Desugar(), tracer)
	hlthReporter := grpcHealth.RegisterHealthChecks(rpc, config.Name, log)

	server := projector.NewPlaceholderProjectorLogic(log, &tracer)
	evented.RegisterProjectorServer(rpc, server)
	hlthReporter.OK()

	lis, err := support.OpenPort(config.Port, log)

	err = rpc.Serve(lis)
	if err != nil {
		log.Error(err)
	}
}
