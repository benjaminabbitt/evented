package memory

import (
	"context"
	"github.com/benjaminabbitt/evented/repository/snapshots"
	"github.com/benjaminabbitt/evented/support"
	"github.com/cucumber/godog"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

var log *zap.SugaredLogger
var sut snapshots.SnapshotStorer

func InitializeTestSuite(ctx *godog.TestSuiteContext) {
	log = support.Log()
}

func InitializeScenario(s *godog.ScenarioContext) {
	sut, _ = NewSnapshotRepoMemory(log)
	s.Step(`^I should be able to retrieve a snapshot with id ([0-9a-fA-F]{8}\-[0-9a-fA-F]{4}\-[0-9a-fA-F]{4}\-[0-9a-fA-F]{4}\-[0-9a-fA-F]{12}) and sequence (\d+)$`, iShouldBeAbleToRetrieveASnapshotWithIdAndSequence)
	s.Step(`^I store a snapshot with id ([0-9a-fA-F]{8}\-[0-9a-fA-F]{4}\-[0-9a-fA-F]{4}\-[0-9a-fA-F]{4}\-[0-9a-fA-F]{12}) and sequence (\d+)$`, iStoreASnapshotWithId)
}

func iShouldBeAbleToRetrieveASnapshotWithIdAndSequence(id string) error {
	uid, _ := uuid.Parse(id)
	_, _ = sut.Get(context.Background(), uid)
	return nil
}

func iStoreASnapshotWithId(id string) error {
	uid, _ := uuid.Parse(id)
	_ = sut.Put(context.Background(), uid, nil)
	return nil
}
