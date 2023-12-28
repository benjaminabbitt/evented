package main

import (
	"context"
	"github.com/benjaminabbitt/evented/applications/command/query-handler/configuration"
	"github.com/benjaminabbitt/evented/applications/command/query-handler/eventQueryServer"
	"github.com/benjaminabbitt/evented/proto/gen/github.com/benjaminabbitt/evented/proto/evented"
	"github.com/benjaminabbitt/evented/repository/events/mongo"
	"github.com/benjaminabbitt/evented/support"
	"github.com/benjaminabbitt/evented/support/grpcHealth"
	"github.com/benjaminabbitt/evented/support/grpcWithInterceptors"
	"github.com/benjaminabbitt/evented/support/jaeger"
)

func main() {
	log := support.Log()
	support.LogStartup(log, "Evented Query Handler")
	defer func() {
		if err := log.Sync(); err != nil {
			log.Errorw("Error syncing logs", err)
		}
	}()
	config := &configuration.Configuration{}
	config = support.Initialize(log, config).(*configuration.Configuration)

	mongoUrl := config.EventStore.Mongodb.Url
	databaseName := config.EventStore.Mongodb.Name
	collectionName := config.EventStore.Mongodb.Collection

	repo, err := mongo.NewEventRepoMongo(context.Background(), mongoUrl, databaseName, collectionName, log)
	if err != nil {
		log.Error(err)
	}

	tracer, closer := jaeger.SetupJaeger(config.Name, log)
	defer jaeger.CloseJaeger(closer, log)

	rpc := grpcWithInterceptors.GenerateConfiguredServer(log.Desugar(), tracer)
	hlthReporter := grpcHealth.RegisterHealthChecks(rpc, config.Name, log)
	server := eventQueryServer.NewEventQueryServer(config.TargetSize, repo, log)
	evented.RegisterEventQueryServer(rpc, server)
	hlthReporter.OK()

	lis, err := support.OpenPort(config.Port, log)

	err = rpc.Serve(lis)
	if err != nil {
		log.Error(err)
	}
}
