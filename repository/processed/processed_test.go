package processed

import (
	"context"
	"fmt"
	"github.com/benjaminabbitt/evented/support"
	"github.com/benjaminabbitt/evented/support/dockerTestSuite"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"testing"
)

func TestLastProcessedRepoSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping docker assisted integration tests in short mode")
	} else {
		suite.Run(t, new(LastProcessedRepo))
	}
}

type LastProcessedRepo struct {
	suite.Suite
	log             *zap.SugaredLogger
	systemUnderTest Processed
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
	dait.CreateNewContainer("mongo", []uint16{27017})
	defer dait.StopContainer()
	processedRepo := NewProcessedClient(
		fmt.Sprintf("mongodb://localhost:%d", dait.Ports[0].PublicPort),
		"databaseName",
		"collectionName",
		o.log,
	)
	ctx := context.Background()
	_ = processedRepo.Received(ctx, o.sampleId, 0)
	last, _ := processedRepo.LastReceived(ctx, o.sampleId)
	o.Assert().Equal(uint32(0), last)
}

func (o *LastProcessedRepo) Test_SequenceGreaterThanZero() {
	dait := dockerTestSuite.DockerAssistedIntegrationTest{}
	dait.CreateNewContainer("mongo", []uint16{27017})
	defer dait.StopContainer()
	processedRepo := NewProcessedClient(
		fmt.Sprintf("mongodb://localhost:%d", dait.Ports[0].PublicPort),
		"databaseName",
		"collectionName",
		o.log,
	)
	ctx := context.Background()
	_ = processedRepo.Received(ctx, o.sampleId, 0)
	_ = processedRepo.Received(ctx, o.sampleId, 1)
	last, _ := processedRepo.LastReceived(ctx, o.sampleId)
	o.Assert().Equal(uint32(1), last)
}
