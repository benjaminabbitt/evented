package main

import (
	"github.com/benjaminabbitt/evented/applications/eventQueryHandler/configuration"
	"github.com/benjaminabbitt/evented/repository/events/mongo"
	"github.com/benjaminabbitt/evented/support"
)

func main() {
	log := support.Log()
	defer log.Sync()

	config := configuration.Configuration{}
	config.Initialize("eventQueryHandlerTest", log)

	mongoUrl := config.DatabaseURL()
	databaseName := config.DatabaseName()
	collectionName := config.DatabaseCollection()

	repo, err := mongo.NewEventRepoMongo(mongoUrl, databaseName, collectionName, log)
	if err != nil {
		log.Error(err)
	}
	server := NewEventQueryServer(config.EventBookTargetSize(), repo, log)

	port := config.Port()
	log.Infow("Starting Business Server...", "port", port)
	err = server.Listen(port)
	if err != nil {
		log.Error(err)
	}
}
