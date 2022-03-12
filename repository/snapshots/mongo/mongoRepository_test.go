package mongo

import (
	"context"
	"fmt"
	"github.com/benjaminabbitt/evented/repository/snapshots"
	"github.com/benjaminabbitt/evented/support"
	"github.com/benjaminabbitt/evented/support/dockerTestSuite"
	"github.com/cucumber/godog"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type MongoSnapshotRepositorySuite struct {
	log  *zap.SugaredLogger
	dait *dockerTestSuite.DockerAssistedIntegrationTest
	sut  snapshots.SnapshotStorer
}

func (suite *MongoSnapshotRepositorySuite) InitializeTestSuite(ctx *godog.TestSuiteContext) {
	suite.log = support.Log()
}

func (suite *MongoSnapshotRepositorySuite) InitializeScenario(s *godog.ScenarioContext) {
	suite.dait = &dockerTestSuite.DockerAssistedIntegrationTest{}
	err := suite.dait.CreateNewContainer("mongo", []uint16{27017})
	if err != nil {
		suite.log.Error(err)
	}
	suite.sut = NewSnapshotMongoRepo(fmt.Sprintf("mongodb://localhost:%d", suite.dait.PublicPort()), "test", suite.log)
	s.Step(`^I should be able to retrieve a snapshot with id ([0-9a-fA-F]{8}\-[0-9a-fA-F]{4}\-[0-9a-fA-F]{4}\-[0-9a-fA-F]{4}\-[0-9a-fA-F]{12}) and sequence (\d+)$`, suite.iShouldBeAbleToRetrieveASnapshotWithIdAndSequence)
	s.Step(`^I store a snapshot with id ([0-9a-fA-F]{8}\-[0-9a-fA-F]{4}\-[0-9a-fA-F]{4}\-[0-9a-fA-F]{4}\-[0-9a-fA-F]{12}) and sequence (\d+)$`, suite.iStoreASnapshotWithId)
}

func (suite *MongoSnapshotRepositorySuite) iShouldBeAbleToRetrieveASnapshotWithIdAndSequence(id string) error {
	uid, _ := uuid.Parse(id)
	_, _ = suite.sut.Get(context.Background(), uid)
	return nil
}

func (suite *MongoSnapshotRepositorySuite) iStoreASnapshotWithId(id string) error {
	uid, _ := uuid.Parse(id)
	_ = suite.sut.Put(context.Background(), uid, nil)
	return nil
}
