package main

import (
	"github.com/benjaminabbitt/evented/repository/events/mongo"
	"github.com/benjaminabbitt/evented/support"
)

func main() {
	log := support.Log()
	defer log.Sync()

	config := Configuration{}
	config.Initialize(log)

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
	server.Listen(port)
}
