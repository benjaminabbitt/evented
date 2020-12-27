package memory

import (
	"context"
	"fmt"
	evented_core "github.com/benjaminabbitt/evented/proto/evented/core"
	"github.com/benjaminabbitt/evented/repository/events/memory"
	"github.com/benjaminabbitt/evented/support"
	"github.com/benjaminabbitt/evented/support/dockerTestSuite"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"testing"
)

func TestLastProcessedRepoSuite(t *testing.T) {
	if !testing.Short() {
		suite.Run(t, new(LastProcessedRepo))
	}
}

type LastProcessedRepo struct {
	suite.Suite
	log             *zap.SugaredLogger
	systemUnderTest memory.EventRepoMemory
	sampleId        uuid.UUID
	mongoId         [12]byte
}

func (o *LastProcessedRepo) SetupTest() {
	o.log = support.Log()
	defer o.log.Sync()

	id, _ := uuid.Parse("c5c10714-2272-4329-809c-38344e318279")
	o.sampleId = id
	o.systemUnderTest = memory.MakeEventRepoMemory(o.log)
}

func (o *LastProcessedRepo) Test_SequenceZero() {
	dait := dockerTestSuite.DockerAssistedIntegrationTest{}
	_ = dait.CreateNewContainer("mongo", []uint16{27017})
	defer dait.StopContainer()
	ctx := context.Background()

	snapshot := &evented_core.Snapshot{
		Sequence: 0,
		State:    nil,
	}

	_ = o.systemUnderTest.Add(ctx, o.sampleId, snapshot)
	last, _ := repo.Get(ctx, o.sampleId)
	o.log.Info(last.Sequence)
}

func (o *LastProcessedRepo) Test_SequenceGreaterThanZero() {
	dait := dockerTestSuite.DockerAssistedIntegrationTest{}
	_ = dait.CreateNewContainer("mongo", []uint16{27017})
	defer dait.StopContainer()
	repo := NewSnapshotMongoRepo(
		fmt.Sprintf("mongodb://localhost:%d", dait.Ports[0].PublicPort),
		"test",
		o.log,
	)

	zero := &evented_core.Snapshot{
		Sequence: 0,
		State:    nil,
	}

	one := &evented_core.Snapshot{
		Sequence: 1,
		State:    nil,
	}

	ctx := context.Background()
	_ = repo.Put(ctx, o.sampleId, zero)
	_ = repo.Put(ctx, o.sampleId, one)
	last, _ := repo.Get(ctx, o.sampleId)
	o.log.Info(last.Sequence)
}
