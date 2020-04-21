package main

import (
	"fmt"
	"github.com/benjaminabbitt/evented/applications/commandHandler/business/client"
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
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var log *zap.SugaredLogger

func main() {
	var name *string = flag.String("appName", "", "The name of the application.  This is used in a number of places, from configuration file name, to queue names.")
	var configPath *string = flag.String("configPath", ".", "The configuration path of the application.  Full config will be located at $configpath/$appName.yaml")
	flag.Parse()

	err := support.SetupConfig(name, configPath, flag.CommandLine)
	if err != nil {
		log.Error(err)
	}
	log = support.Log()
	defer log.Sync()

	businessAddress := viper.GetString("business.address")
	commandHandlerPort := uint16(viper.GetUint("port"))
	log.Infow("Starting Command Handler", "port", commandHandlerPort)
	businessClient, _ := client.NewBusinessClient(businessAddress, log)

	eventRepo, err := setupEventRepo(log)
	ssRepo := setupSnapshotRepo()
	domain := viper.GetString("domain")

	repo := eventBook.RepositoryBasic{
		EventRepo:    eventRepo,
		SnapshotRepo: ssRepo,
		Domain:       domain,
	}

	handlers := transport.NewTransportHolder(log)

	sagaConfig := viper.GetStringMap("sync.sagas")
	for name, _ := range sagaConfig {
		url := viper.GetString("sync.sagas." + name + ".url")
		sagaConn := grpcWithInterceptors.GenerateConfiguredConn(url, log)
		log.Infow("Synchronous Saga Configuration:", name, url)
		handlers.Add(saga.NewGRPCSagaClient(sagaConn))
	}

	projectorConfig := viper.GetStringSlice("sync.projectors")
	for _, url := range projectorConfig {
		log.Infow("Synchronous Projector Configuration:", "host", url)
		projectorConn := grpcWithInterceptors.GenerateConfiguredConn(url, log)
		handlers.Add(projector.NewGRPCProjector(projectorConn))
	}

	handlers.Add(setupServiceBus(domain))

	server := framework.NewServer(
		repo,
		*handlers,
		businessClient,
		log,
	)
	server.Listen(commandHandlerPort)
}

func setupSnapshotRepo() (repo snapshots.SnapshotStorer) {
	configurationKey := "snapshotStore"
	typee := viper.GetString("snapshotStore.type")
	mongodb := "mongodb"
	if typee == mongodb {
		url := viper.GetString(fmt.Sprintf("%s.%s.url", configurationKey, mongodb))
		dbName := viper.GetString(fmt.Sprintf("%s.%s.database", configurationKey, mongodb))
		repo = snapshot_mongo.NewSnapshotMongoRepo(url, dbName, log)
	}
	return repo
}

func setupServiceBus(domain string) (trans chan *evented_core.EventBook) {
	trans = make(chan *evented_core.EventBook)
	configurationKey := "transport"
	amqpText := "amqp"
	typee := viper.GetString(fmt.Sprintf("%s.type", configurationKey))
	if typee == amqpText {
		url := viper.GetString(fmt.Sprintf("%s.%s.url", configurationKey, amqpText))
		exchange := viper.GetString(fmt.Sprintf("%s.%s.exchange", configurationKey, amqpText))

		_ = sender.NewAMQPSender(trans, url, exchange, log)
		return trans
	}
	return nil
}
func setupEventRepo(log *zap.SugaredLogger) (repo events.EventStorer, err error) {
	configurationKey := "eventStore"
	typee := viper.GetString("eventstore.type")
	mongodb := "mongodb"
	if typee == mongodb {
		url := viper.GetString(fmt.Sprintf("%s.%s.url", configurationKey, mongodb))
		dbName := viper.GetString(fmt.Sprintf("%s.%s.database", configurationKey, mongodb))
		collectionName := viper.GetString(fmt.Sprintf("%s.%s.collection", configurationKey, mongodb))
		log.Infow("Using MongoDb for Event Store", "url", url, "dbName", dbName)
		repo, err := event_mongo.NewEventRepoMongo(url, dbName, collectionName, log)
		if err != nil {
			return nil, err
		}
		return repo, nil
	}
	return repo, nil
}
