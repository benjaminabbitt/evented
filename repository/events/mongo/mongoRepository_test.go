package mongo

import (
	"context"
	"fmt"
	core "github.com/benjaminabbitt/evented/proto/evented/core"
	"github.com/benjaminabbitt/evented/repository/events"
	"github.com/benjaminabbitt/evented/support"
	"github.com/benjaminabbitt/evented/support/cucumber"
	"github.com/benjaminabbitt/evented/support/dockerTestSuite"
	"github.com/cucumber/godog"
	"github.com/cucumber/messages-go/v10"
	"github.com/golang/protobuf/ptypes"
	timestamppb "github.com/golang/protobuf/ptypes/timestamp"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"strconv"
	"time"
)

type MongoRepositorySuite struct {
	log    *zap.SugaredLogger
	id     uuid.UUID
	events []*evented.EventPage
	dait   *dockerTestSuite.DockerAssistedIntegrationTest
	sut    events.EventStorer
}

func (suite *MongoRepositorySuite) InitializeTestSuite(ctx *godog.TestSuiteContext) {
	suite.log = support.Log()
}

func (suite *MongoRepositorySuite) InitializeScenario(s *godog.ScenarioContext) {
	suite.dait = &dockerTestSuite.DockerAssistedIntegrationTest{}
	err := suite.dait.CreateNewContainer("mongo", []uint16{27017})
	if err != nil {
		suite.log.Error(err)
	}

	suite.id, _ = uuid.NewRandom()
	suite.sut, _ = NewEventRepoMongo(context.Background(), fmt.Sprintf("mongodb://localhost:%d", suite.dait.Ports[0].PublicPort), "test", "events", suite.log)
	s.Step(`^I should be able to retrieve it by its coordinates:$`, suite.iShouldBeAbleToRetrieveItByItsCoordinates)
	s.Step(`^I store the event:$`, suite.iStoreTheEvent)
	s.Step(`^a populated database:$`, suite.aPopulatedDatabase)
	s.Step(`^I should get these events:$`, suite.iShouldGetTheseEvents)
	s.Step(`^I retrieve a subset of events ending at event (\d+)$`, suite.iRetrieveASubsetOfEventsEndingAtEvent)
	s.Step(`^I retrieve a subset of events from (\d+) to (\d+)$`, suite.iRetrieveASubsetOfEventsFromTo)
	s.Step(`^I retrieve a subset of events starting from value (\d+)$`, suite.iRetrieveASubsetOfEventsStartingFromValue)
	s.Step(`^I retrieve all events$`, suite.iRetrieveAllEvents)
}

func (suite *MongoRepositorySuite) iShouldBeAbleToRetrieveItByItsCoordinates(arg1 *messages.PickleStepArgument_PickleTable) error {
	suite.id, suite.events = suite.extractPickleTableToEvents(arg1)
	ch := make(chan *evented.EventPage)
	_ = suite.sut.Get(context.Background(), ch, suite.id)
	return cucumber.AssertExpectedAndActual(assert.Equal, suite.events[0], <-ch, "", "")
}

func (suite *MongoRepositorySuite) extractPickleTableToEvents(arg *messages.PickleStepArgument_PickleTable) (id uuid.UUID, events []*evented.EventPage) {
	for i, row := range arg.GetRows() {
		if i == 0 { //header
			continue
		}
		var sequence uint32
		var force bool
		var ts *timestamppb.Timestamp
		for j, cell := range row.GetCells() {
			switch j {
			case 0:
				id, _ = uuid.Parse(cell.Value)
			case 1:
				if cell.Value == "force" {
					force = true
				} else {
					sequence, _ = cucumber.Uint64ToUint32WithErrorPassthrough(strconv.ParseUint(cell.Value, 10, 32))
				}
			case 2:
				t, _ := time.Parse(time.RFC3339Nano, cell.Value)
				ts, _ = ptypes.TimestampProto(t)
			}

		}

		event := &evented.EventPage{
			CreatedAt:   ts,
			Event:       nil,
			Synchronous: false,
		}
		if force {
			event.Sequence = &evented.EventPage_Force{}
		} else {
			event.Sequence = &evented.EventPage_Num{Num: sequence}
		}
		events = append(events, event)
	}
	return id, events
}

func (suite *MongoRepositorySuite) iStoreTheEvent(arg1 *messages.PickleStepArgument_PickleTable) error {
	suite.id, suite.events = suite.extractPickleTableToEvents(arg1)
	_ = suite.sut.Add(context.Background(), suite.id, suite.events)
	return nil
}

func (suite *MongoRepositorySuite) aPopulatedDatabase(arg1 *messages.PickleStepArgument_PickleTable) error {
	suite.id, suite.events = suite.extractPickleTableToEvents(arg1)
	_ = suite.sut.Add(context.Background(), suite.id, suite.events)
	return nil
}

func (suite *MongoRepositorySuite) iShouldGetTheseEvents(arg1 *messages.PickleStepArgument_PickleTable) error {
	_, expectedEvents := suite.extractPickleTableToEvents(arg1)
	return cucumber.AssertExpectedAndActual(assert.Equal, expectedEvents, suite.events, "", "")
}

func (suite *MongoRepositorySuite) drainChannel(ch chan *evented.EventPage) (pages []*evented.EventPage) {
	for page := range ch {
		pages = append(pages, page)
	}
	return pages
}

func (suite *MongoRepositorySuite) iRetrieveASubsetOfEventsEndingAtEvent(end int) error {
	ch := make(chan *evented.EventPage)
	_ = suite.sut.GetTo(context.Background(), ch, suite.id, uint32(end))
	suite.events = suite.drainChannel(ch)
	return nil
}

func (suite *MongoRepositorySuite) iRetrieveASubsetOfEventsFromTo(start, end int) error {
	ch := make(chan *evented.EventPage)
	_ = suite.sut.GetFromTo(context.Background(), ch, suite.id, uint32(start), uint32(end))
	suite.events = suite.drainChannel(ch)
	return nil
}

func (suite *MongoRepositorySuite) iRetrieveASubsetOfEventsStartingFromValue(start int) error {
	ch := make(chan *evented.EventPage)
	_ = suite.sut.GetFrom(context.Background(), ch, suite.id, uint32(start))
	suite.events = suite.drainChannel(ch)
	return nil
}

func (suite *MongoRepositorySuite) iRetrieveAllEvents() error {
	ch := make(chan *evented.EventPage)
	_ = suite.sut.Get(context.Background(), ch, suite.id)
	suite.events = suite.drainChannel(ch)
	return nil
}
