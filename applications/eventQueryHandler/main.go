package main

import (
	"context"
	"github.com/benjaminabbitt/evented/applications/eventQueryHandler/configuration"
	"github.com/benjaminabbitt/evented/applications/eventQueryHandler/eventQueryServer"
	"github.com/benjaminabbitt/evented/repository/events/mongo"
	"github.com/benjaminabbitt/evented/support"
	"github.com/benjaminabbitt/evented/support/jaeger"
)

func main() {
	log := support.Log()
	defer log.Sync()

	config := configuration.Configuration{}
	config.Initialize(log)

	mongoUrl := config.DatabaseURL()
	databaseName := config.DatabaseName()
	collectionName := config.DatabaseCollection()

	repo, err := mongo.NewEventRepoMongo(context.Background(), mongoUrl, databaseName, collectionName, log)
	if err != nil {
		log.Error(err)
	}
	server := eventQueryServer.NewEventQueryServer(config.EventBookTargetSize(), repo, log)

	tracer, closer := jaeger.SetupJaeger(config.AppName(), log)
	defer closer.Close()

	port := config.Port()
	log.Infow("Starting Business Server...", "port", port)
	err = server.Listen(port, tracer)
	if err != nil {
		log.Error(err)
	}
}
