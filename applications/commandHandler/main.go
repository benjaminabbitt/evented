package main

import (
	"fmt"
	"github.com/benjaminabbitt/evented"
	"github.com/benjaminabbitt/evented/applications/commandHandler/business/client"
	"github.com/benjaminabbitt/evented/applications/commandHandler/framework"
	"github.com/benjaminabbitt/evented/applications/commandHandler/transport"
	"github.com/benjaminabbitt/evented/applications/commandHandler/transport/async/amqp"
	mockProjector "github.com/benjaminabbitt/evented/applications/commandHandler/transport/sync/projector/mock"
	mockSaga "github.com/benjaminabbitt/evented/applications/commandHandler/transport/sync/saga/mock"
	"github.com/benjaminabbitt/evented/repository/eventBook"
	"github.com/benjaminabbitt/evented/repository/events"
	memoryRepository "github.com/benjaminabbitt/evented/repository/events/event-memory"
	"github.com/benjaminabbitt/evented/repository/events/mongo"
	"github.com/benjaminabbitt/evented/repository/snapshots"
	snapshot_memory "github.com/benjaminabbitt/evented/repository/snapshots/snapshot-memory"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

const (
	NAME = "commandHandler"
)

var log *zap.SugaredLogger
var errh *evented.ErrLogger

func main() {
	setupConfig()

	logger, _ := zap.NewDevelopment(zap.AddCaller())
	log = logger.Sugar()

	defer log.Sync()
	log.Infow("Logger Configured")

	errh = &evented.ErrLogger{log}

	log.Infow("Configuration Read")

	businessAddress := viper.GetString("business.address")
	commandHandlerPort := uint16(viper.GetUint("port"))
	log.Infow("Starting Command Handler", "port", commandHandlerPort)
	businessClient, _ := client.NewBusinessClient(businessAddress, log)

	eventRepo := setupEventRepo()
	ssRepo := setupSnapshotRepo()
	domain := viper.GetString("domain")

	repo := eventBook.Repository{
		EventRepo:    eventRepo,
		SnapshotRepo: ssRepo,
		Domain:       domain,
	}

	handlers := transport.NewTransportHolder(log)

	handlers.Add(mockSaga.NewSagaClient(log))
	handlers.Add(mockProjector.NewProjectorClient(log))
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

func setupConfig() {
	viper.SetConfigName(NAME)
	viper.SetConfigType("yaml")

	viper.AddConfigPath(".")
	viper.AddConfigPath("c:/temp/")

	viper.SetEnvPrefix(NAME)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			log.Warn(err)
		} else {
			log.Fatal(err)
		}
	}
}

func setupEventRepo()(repo events.EventRepository){
	log.Infow("test")
	configurationKey := "eventStore"
	typee := viper.GetString("eventstore.type")
	mongodb := "mongodb"
	memory := "memory"
	if typee == mongodb {
		url := viper.GetString(fmt.Sprintf("%s.%s.url", configurationKey, mongodb))
		dbName := viper.GetString(fmt.Sprintf("%s.%s.database", configurationKey, mongodb))
		log.Infow("Using MongoDb for Event Store", "url", url, "dbName", dbName)
		repo = mongo.NewMongoClient(url, dbName, log, errh)
	}else if typee == memory{
		repo = memoryRepository.NewMemoryRepository()
	}
	return repo
}

func setupSnapshotRepo()(repo snapshots.SnapshotRepo){
	configurationKey := "snapshotStore"
	typee := viper.GetString("snapshotStore.type")
	mongodb := "mongodb"
	memory := "memory"
	if typee == mongodb{
		url := viper.GetString(fmt.Sprintf("%s.%s.url", configurationKey, mongodb))
		dbName := viper.GetString(fmt.Sprintf("%s.%s.database", configurationKey, mongodb))
		log.Warnw("Read Mongo preference for snapshot repo.  Snapshot Mongo incomplete.  Configuring wtih Memory for now...", "url", url, "dbName",dbName)
		repo = snapshot_memory.NewSSMemoryRepository()
	}else if typee == memory{
		repo = snapshot_memory.NewSSMemoryRepository()
	}
	return repo
}

func setupServiceBus(domain string)(transport transport.Transport){
	configurationKey := "transport"
	amqpText := "amqp"
	typee := viper.GetString(fmt.Sprintf("%s.type", configurationKey))
	if typee == amqpText {
		url := viper.GetString(fmt.Sprintf("%s.%s.url", configurationKey, amqpText))
		exchange := viper.GetString(fmt.Sprintf("%s.%s.exchange", configurationKey, amqpText))
		client := amqp.NewAMQPClient(url, "evented_"+exchange, log, errh)
		return client
	}
	return nil
}