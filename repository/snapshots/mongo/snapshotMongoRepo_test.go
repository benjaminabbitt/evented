package mongo

import (
	"context"
	"fmt"
	evented_core "github.com/benjaminabbitt/evented/proto/core"
	"github.com/benjaminabbitt/evented/support"
	"github.com/benjaminabbitt/evented/support/dockerTestSuite"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"testing"
)

func TestLastProcessedRepoSuite(t *testing.T) {
	suite.Run(t, new(LastProcessedRepo))
}

type LastProcessedRepo struct {
	suite.Suite
	log             *zap.SugaredLogger
	systemUnderTest SnapshotMongoRepo
	sampleId        uuid.UUID
	mongoId         [12]byte
}

func (o *LastProcessedRepo) SetupTest() {
	o.log = support.Log()
	defer o.log.Sync()

	id, _ := uuid.Parse("c5c10714-2272-4329-809c-38344e318279")
	o.sampleId = id
	o.mongoId = [12]byte{197, 193, 7, 20, 34, 114, 67, 41, 128, 156, 56, 52}
}

func (o *LastProcessedRepo) Test_SequenceZero() {
	dait := dockerTestSuite.DockerAssistedIntegrationTest{}
	_ = dait.CreateNewContainer("mongo", []uint16{27017})
	defer dait.StopContainer()
	repo := NewSnapshotMongoRepo(
		fmt.Sprintf("mongodb://localhost:%d", dait.Ports[0].PublicPort),
		"test",
		o.log,
	)
	ctx := context.Background()

	snapshot := &evented_core.Snapshot{
		Sequence: 0,
		State:    nil,
	}

	_ = repo.Put(ctx, o.sampleId, snapshot)
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
