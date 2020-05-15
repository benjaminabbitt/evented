package main

import (
	"github.com/benjaminabbitt/evented/applications/commandHandler/business/client"
	"github.com/benjaminabbitt/evented/applications/commandHandler/configuration"
	"github.com/benjaminabbitt/evented/applications/commandHandler/framework"
	"github.com/benjaminabbitt/evented/applications/commandHandler/framework/transport"
	eventedcore "github.com/benjaminabbitt/evented/proto/evented/core"
	"github.com/benjaminabbitt/evented/repository/eventBook"
	"github.com/benjaminabbitt/evented/repository/events"
	eventmongo "github.com/benjaminabbitt/evented/repository/events/mongo"
	"github.com/benjaminabbitt/evented/repository/snapshots"
	snapshotmongo "github.com/benjaminabbitt/evented/repository/snapshots/mongo"
	"github.com/benjaminabbitt/evented/support"
	"github.com/benjaminabbitt/evented/support/grpcWithInterceptors"
	"github.com/benjaminabbitt/evented/transport/async/amqp/sender"
	"github.com/benjaminabbitt/evented/transport/sync/projector"
	"github.com/benjaminabbitt/evented/transport/sync/saga"
	"go.uber.org/zap"
)

var log *zap.SugaredLogger

func main() {
	log = support.Log()

	config := configuration.Configuration{}
	config.Initialize("commandHandler", log)

	businessAddress := config.BusinessURL()
	commandHandlerPort := config.Port()
	log.Infow("Starting Command Handler", "port", commandHandlerPort)
	businessClient, _ := client.NewBusinessClient(businessAddress, log)
	log.Infow("Command Handler Started", "port", commandHandlerPort)

	eventRepo, _ := setupEventRepo(config, log)
	ssRepo := setupSnapshotRepo(config)

	repo := eventBook.MakeRepositoryBasic(eventRepo, ssRepo, config.Domain(), log)

	handlers := transport.NewTransportHolder(log)

	for _, url := range config.SagaURLs() {
		log.Infow("Connecting with Saga... ", "url", url)
		sagaConn := grpcWithInterceptors.GenerateConfiguredConn(url, log)
		handlers.Add(saga.NewGRPCSagaClient(sagaConn))
		log.Infow("Connection with Saga Successful", "url", url)
	}

	for _, url := range config.ProjectorURLs() {
		log.Infow("Connecting with Projector... ", "url", url)
		projectorConn := grpcWithInterceptors.GenerateConfiguredConn(url, log)
		handlers.Add(projector.NewGRPCProjector(projectorConn))
		log.Infow("Connection with Projector Successful.", "url", url)
	}

	handlers.Add(setupServiceBus(config))

	server := framework.NewServer(
		repo,
		handlers,
		businessClient,
		log,
	)
	server.Listen(config.Port())
}

func setupSnapshotRepo(config configuration.Configuration) (repo snapshots.SnapshotStorer) {
	return snapshotmongo.NewSnapshotMongoRepo(config.SnapshotStoreURL(), config.SnapshotStoreDatabaseName(), log)
}

func setupServiceBus(config configuration.Configuration) (ch chan *eventedcore.EventBook) {
	ch = make(chan *eventedcore.EventBook)
	trans := sender.NewAMQPSender(ch, config.TransportURL(), config.TransportExchange(), log)
	err := trans.Connect()
	if err != nil {
		log.Error(err)
	}
	trans.Run()
	return ch
}

func setupEventRepo(config configuration.Configuration, log *zap.SugaredLogger) (repo events.EventStorer, err error) {
	repo, err = eventmongo.NewEventRepoMongo(config.EventStoreURL(), config.EventStoreDatabaseName(), config.EventStoreCollectionName(), log)
	if err != nil {
		return nil, err
	}
	return repo, nil
}
