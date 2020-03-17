package main

import (
	"fmt"
	"github.com/benjaminabbitt/evented/business"
	"github.com/benjaminabbitt/evented/framework"
	"github.com/benjaminabbitt/evented/repository/eventBook"
	memoryRepository "github.com/benjaminabbitt/evented/repository/events/event-memory"
	snapshot_memory "github.com/benjaminabbitt/evented/repository/snapshots/snapshot-memory"
	"github.com/benjaminabbitt/evented/transport"
	mockProjector "github.com/benjaminabbitt/evented/transport/projector/mock"
	mockSaga "github.com/benjaminabbitt/evented/transport/saga/mock"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	NAME = "test"
)

func main() {
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

	log.Print(viper.Get("eventstore.type"))
	log.Print(viper.Get(fmt.Sprintf("eventstore.%s.url", viper.GetString("eventstore.Type"))))

	businessAddress := viper.GetString("business.address")
	log.WithFields(log.Fields{
		"address": businessAddress,
	}).Info("Starting Command Handler")
	framework.NewServer(
		eventBook.Repository{EventRepo: memoryRepository.NewMemoryRepository(), SnapshotRepo: snapshot_memory.NewSSMemoryRepository(), Domain: "domain"},
		[]transport.SyncSaga{mockSaga.SagaClient{}},
		[]transport.SyncProjection{mockProjector.NewProjectorClient()},
		[]transport.Saga{mockSaga.SagaClient{}},
		[]transport.Projection{mockProjector.NewProjectorClient()},
		&business.MockBusinessLogicClient{},
	)
	log.WithFields(log.Fields{
		"address": businessAddress,
	}).Info("Command Handler started")
}
