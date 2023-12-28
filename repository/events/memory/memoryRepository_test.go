package memory

import (
	"context"
	"github.com/benjaminabbitt/evented/proto/gen/github.com/benjaminabbitt/evented/proto/evented"
	"github.com/benjaminabbitt/evented/support/cucumber"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"
	"strconv"
	"time"

	"github.com/benjaminabbitt/evented/support"
	"github.com/cucumber/godog"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type MemoryRepositorySuite struct {
	log    *zap.SugaredLogger
	sut    EventRepoMemory
	id     uuid.UUID
	events []*evented.EventPage
}

func (suite *MemoryRepositorySuite) InitializeTestSuite(ctx *godog.TestSuiteContext) {
	suite.log = support.Log()
}

func (suite *MemoryRepositorySuite) InitializeScenario(s *godog.ScenarioContext) {
	suite.sut, _ = NewEventRepoMemory(suite.log)
	s.Step(`^I should be able to retrieve it by its coordinates:$`, suite.iShouldBeAbleToRetrieveItByItsCoordinates)
	s.Step(`^I store the event:$`, suite.iStoreTheEvent)
	s.Step(`^a populated database:$`, suite.aPopulatedDatabase)
	s.Step(`^I should get these events:$`, suite.iShouldGetTheseEvents)
	s.Step(`^I retrieve a subset of events ending at event (\d+)$`, suite.iRetrieveASubsetOfEventsEndingAtEvent)
	s.Step(`^I retrieve a subset of events from (\d+) to (\d+)$`, suite.iRetrieveASubsetOfEventsFromTo)
	s.Step(`^I retrieve a subset of events starting from value (\d+)$`, suite.iRetrieveASubsetOfEventsStartingFromValue)
	s.Step(`^I retrieve all events$`, suite.iRetrieveAllEvents)
}

func (suite *MemoryRepositorySuite) iShouldBeAbleToRetrieveItByItsCoordinates(arg1 *godog.Table) error {
	id, events := suite.extractPickleTableToEvents(arg1)
	ch := make(chan *evented.EventPage)
	_ = suite.sut.Get(context.Background(), ch, id)
	return cucumber.AssertExpectedAndActual(assert.Equal, events[0], <-ch, "", "")
}

func (suite *MemoryRepositorySuite) extractPickleTableToEvents(arg *godog.Table) (id uuid.UUID, events []*evented.EventPage) {
	for i, row := range arg.Rows {
		if i == 0 { //header
			continue
		}
		var sequence uint32
		var force bool
		var ts *timestamppb.Timestamp
		for j, cell := range row.Cells {
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
				ts = timestamppb.New(t)
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

func (suite *MemoryRepositorySuite) iStoreTheEvent(arg1 *godog.Table) error {
	suite.id, suite.events = suite.extractPickleTableToEvents(arg1)
	_ = suite.sut.Add(context.Background(), suite.id, suite.events)
	return nil
}

func (suite *MemoryRepositorySuite) aPopulatedDatabase(arg1 *godog.Table) error {
	suite.id, suite.events = suite.extractPickleTableToEvents(arg1)
	_ = suite.sut.Add(context.Background(), suite.id, suite.events)
	return nil
}

func (suite *MemoryRepositorySuite) iShouldGetTheseEvents(arg1 *godog.Table) error {
	_, expectedEvents := suite.extractPickleTableToEvents(arg1)
	return cucumber.AssertExpectedAndActual(assert.Equal, expectedEvents, suite.events, "", "")
}

func (suite *MemoryRepositorySuite) drainChannel(ch chan *evented.EventPage) (pages []*evented.EventPage) {
	for page := range ch {
		pages = append(pages, page)
	}
	return pages
}

func (suite *MemoryRepositorySuite) iRetrieveASubsetOfEventsEndingAtEvent(end int) error {
	ch := make(chan *evented.EventPage)
	_ = suite.sut.GetTo(context.Background(), ch, suite.id, uint32(end))
	suite.events = suite.drainChannel(ch)
	return nil
}

func (suite *MemoryRepositorySuite) iRetrieveASubsetOfEventsFromTo(start, end int) error {
	ch := make(chan *evented.EventPage)
	_ = suite.sut.GetFromTo(context.Background(), ch, suite.id, uint32(start), uint32(end))
	suite.events = suite.drainChannel(ch)
	return nil
}

func (suite *MemoryRepositorySuite) iRetrieveASubsetOfEventsStartingFromValue(start int) error {
	ch := make(chan *evented.EventPage)
	_ = suite.sut.GetFrom(context.Background(), ch, suite.id, uint32(start))
	suite.events = suite.drainChannel(ch)
	return nil
}

func (suite *MemoryRepositorySuite) iRetrieveAllEvents() error {
	ch := make(chan *evented.EventPage)
	_ = suite.sut.Get(context.Background(), ch, suite.id)
	suite.events = suite.drainChannel(ch)
	return nil
}
