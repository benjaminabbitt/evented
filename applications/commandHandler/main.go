package main

import (
	"github.com/benjaminabbitt/evented/applications/commandHandler/business/client"
	"github.com/benjaminabbitt/evented/applications/commandHandler/configuration"
	"github.com/benjaminabbitt/evented/applications/commandHandler/framework"
	"github.com/benjaminabbitt/evented/applications/commandHandler/framework/transport"
	evented_core "github.com/benjaminabbitt/evented/proto/core"
	"github.com/benjaminabbitt/evented/repository/eventBook"
	"github.com/benjaminabbitt/evented/repository/events"
	event_mongo "github.com/benjaminabbitt/evented/repository/events/mongo"
	"github.com/benjaminabbitt/evented/repository/snapshots"
	snapshot_mongo "github.com/benjaminabbitt/evented/repository/snapshots/mongo"
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

	repo := eventBook.RepositoryBasic{
		EventRepo:    eventRepo,
		SnapshotRepo: ssRepo,
		Domain:       config.Domain(),
	}

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
		*handlers,
		businessClient,
		log,
	)
	server.Listen(config.Port())
}

func setupSnapshotRepo(config configuration.Configuration) (repo snapshots.SnapshotStorer) {
	return snapshot_mongo.NewSnapshotMongoRepo(config.SnapshotStoreURL(), config.SnapshotStoreDatabaseName(), log)
}

func setupServiceBus(config configuration.Configuration) (ch chan *evented_core.EventBook) {
	ch = make(chan *evented_core.EventBook)
	trans := sender.NewAMQPSender(ch, config.TransportURL(), config.TransportExchange(), log)
	trans.Run()
	return ch
}

func setupEventRepo(config configuration.Configuration, log *zap.SugaredLogger) (repo events.EventStorer, err error) {
	repo, err = event_mongo.NewEventRepoMongo(config.EventStoreURL(), config.EventStoreDatabaseName(), config.EventStoreCollectionName(), log)
	if err != nil {
		return nil, err
	}
	return repo, nil
}
