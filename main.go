package main

import (
	"github.com/benjaminabbitt/evented/business"
	"github.com/benjaminabbitt/evented/framework"
	"github.com/benjaminabbitt/evented/repository/eventBook"
	memoryRepository "github.com/benjaminabbitt/evented/repository/events/event-memory"
	snapshot_memory "github.com/benjaminabbitt/evented/repository/snapshots/snapshot-memory"
	"github.com/benjaminabbitt/evented/transport"
	mockProjector "github.com/benjaminabbitt/evented/transport/projector/mock"
	mockSaga "github.com/benjaminabbitt/evented/transport/saga/mock"
	log "github.com/sirupsen/logrus"
)

func main() {
	businessAddress := "localhost:8081"
	log.WithFields(log.Fields{
		"address": businessAddress,
	}).Info("Starting Command Handler")
	framework.NewServer(
		eventBook.Repository{EventRepo: memoryRepository.NewMemoryRepository(), SnapshotRepo: snapshot_memory.NewSSMemoryRepository(), Domain: "domain"},
		[]transport.Saga{mockSaga.SagaClient{}},
		[]transport.Projection{mockProjector.NewProjectorClient()},
		[]transport.Saga{mockSaga.SagaClient{}},
		[]transport.Projection{mockProjector.NewProjectorClient()},
		&business.MockBusinessLogicClient{},
	)
	log.WithFields(log.Fields{
		"address": businessAddress,
	}).Info("Command Handler started")
}
