package main

import (
	"fmt"
	"github.com/benjaminabbitt/evented"
	"github.com/benjaminabbitt/evented/applications/commandHandler/business/client"
	"github.com/benjaminabbitt/evented/applications/commandHandler/framework"
	"github.com/benjaminabbitt/evented/applications/commandHandler/framework/transport"
	"github.com/benjaminabbitt/evented/repository/eventBook"
	"github.com/benjaminabbitt/evented/repository/events"
	event_mongo "github.com/benjaminabbitt/evented/repository/events/mongo"
	"github.com/benjaminabbitt/evented/repository/snapshots"
	snapshot_mongo "github.com/benjaminabbitt/evented/repository/snapshots/mongo"
	snapshot_memory "github.com/benjaminabbitt/evented/repository/snapshots/snapshot-memory"
	"github.com/benjaminabbitt/evented/support"
	"github.com/benjaminabbitt/evented/transport/async"
	"github.com/benjaminabbitt/evented/transport/async/amqp"
	"github.com/benjaminabbitt/evented/transport/sync/projector"
	"github.com/benjaminabbitt/evented/transport/sync/saga"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

const (
	NAME = "commandHandler"
)

var log *zap.SugaredLogger
var errh *evented.ErrLogger

func main() {
	var name *string = flag.String("appName", "", "The name of the application.  This is used in a number of places, from configuration file name, to queue names.")
	var configPath *string = flag.String("configPath", ".", "The configuration path of the application.  Full config will be located at $configpath/$appName.yaml")
	flag.Parse()

	err := support.SetupConfig(name, configPath, flag.CommandLine)
	log, errh = support.Log()
	defer log.Sync()
	errh.LogIfErr(err, "Error configuring application.")

	businessAddress := viper.GetString("business.address")
	commandHandlerPort := uint16(viper.GetUint("port"))
	log.Infow("Starting Command Handler", "port", commandHandlerPort)
	businessClient, _ := client.NewBusinessClient(businessAddress, log)

	eventRepo, err := setupEventRepo(log, errh)
	errh.LogIfErr(err, "Error configuring event repo")
	ssRepo := setupSnapshotRepo()
	domain := viper.GetString("domain")

	repo := eventBook.RepositoryBasic{
		EventRepo:    *eventRepo,
		SnapshotRepo: ssRepo,
		Domain:       domain,
	}

	handlers := transport.NewTransportHolder(log)

	handlers.Add(saga.MockSagaClient{})
	handlers.Add(projector.MockProjectorClient{})
	handlers.Add(setupServiceBus(domain))

	server := framework.NewServer(
		repo,
		*handlers,
		businessClient,
		log,
		errh,
	)
	server.Listen(commandHandlerPort)
}

func setupSnapshotRepo() (repo snapshots.SnapshotRepo) {
	configurationKey := "snapshotStore"
	typee := viper.GetString("snapshotStore.type")
	mongodb := "mongodb"
	memory := "memory"
	if typee == mongodb {
		url := viper.GetString(fmt.Sprintf("%s.%s.url", configurationKey, mongodb))
		dbName := viper.GetString(fmt.Sprintf("%s.%s.database", configurationKey, mongodb))
		repo = snapshot_mongo.NewSnapshotMongoRepo(url, dbName, log, errh)
	} else if typee == memory {
		repo = snapshot_memory.NewSSMemoryRepository()
	}
	return repo
}

func setupServiceBus(domain string) (transport async.Transport) {
	configurationKey := "transport"
	amqpText := "amqp"
	typee := viper.GetString(fmt.Sprintf("%s.type", configurationKey))
	if typee == amqpText {
		url := viper.GetString(fmt.Sprintf("%s.%s.url", configurationKey, amqpText))
		exchange := viper.GetString(fmt.Sprintf("%s.%s.exchange", configurationKey, amqpText))
		client := amqp.NewAMQPClient(url, exchange, log, errh)
		return client
	}
	return nil
}
func setupEventRepo(log *zap.SugaredLogger, errh *evented.ErrLogger) (repo *events.EventRepository, err error) {
	configurationKey := "eventStore"
	typee := viper.GetString("eventstore.type")
	mongodb := "mongodb"
	if typee == mongodb {
		url := viper.GetString(fmt.Sprintf("%s.%s.url", configurationKey, mongodb))
		dbName := viper.GetString(fmt.Sprintf("%s.%s.database", configurationKey, mongodb))
		collectionName := viper.GetString(fmt.Sprintf("%s.%s.collection", configurationKey, mongodb))
		log.Infow("Using MongoDb for Event Store", "url", url, "dbName", dbName)
		tempRepo := event_mongo.NewEventRepoMongo(url, dbName, collectionName, log, errh)
		repo = &tempRepo
	}
	return repo, nil
}
