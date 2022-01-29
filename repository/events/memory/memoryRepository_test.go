package memory

import (
	"context"
	"github.com/benjaminabbitt/evented/proto/gen/github.com/benjaminabbitt/evented/proto/evented"

	"github.com/benjaminabbitt/evented/support"
	"github.com/benjaminabbitt/evented/support/cucumber"
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

type MemoryRepositorySuite struct {
	log    *zap.SugaredLogger
	sut    EventRepoMemory
	id     uuid.UUID
	events []*evented.EventPage
}

func (s *MemoryRepositorySuite) InitializeTestSuite(ctx *godog.TestSuiteContext) {
	s.log = support.Log()
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

func (s *MemoryRepositorySuite) iShouldBeAbleToRetrieveItByItsCoordinates(arg1 *messages.PickleStepArgument_PickleTable) error {
	id, events := s.extractPickleTableToEvents(arg1)
	ch := make(chan *evented.EventPage)
	_ = s.sut.Get(context.Background(), ch, id)
	return cucumber.AssertExpectedAndActual(assert.Equal, events[0], <-ch, "", "")
}

func (s *MemoryRepositorySuite) extractPickleTableToEvents(arg *messages.PickleStepArgument_PickleTable) (id uuid.UUID, events []*evented.EventPage) {
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

func (s *MemoryRepositorySuite) iStoreTheEvent(arg1 *messages.PickleStepArgument_PickleTable) error {
	s.id, s.events = s.extractPickleTableToEvents(arg1)
	_ = s.sut.Add(context.Background(), s.id, s.events)
	return nil
}

func (s *MemoryRepositorySuite) aPopulatedDatabase(arg1 *messages.PickleStepArgument_PickleTable) error {
	s.id, s.events = s.extractPickleTableToEvents(arg1)
	_ = s.sut.Add(context.Background(), s.id, s.events)
	return nil
}

func (s *MemoryRepositorySuite) iShouldGetTheseEvents(arg1 *messages.PickleStepArgument_PickleTable) error {
	_, expectedEvents := s.extractPickleTableToEvents(arg1)
	return cucumber.AssertExpectedAndActual(assert.Equal, expectedEvents, s.events, "", "")
}

func (s *MemoryRepositorySuite) drainChannel(ch chan *evented.EventPage) (pages []*evented.EventPage) {
	for page := range ch {
		pages = append(pages, page)
	}
	return pages
}

func (s *MemoryRepositorySuite) iRetrieveASubsetOfEventsEndingAtEvent(end int) error {
	ch := make(chan *evented.EventPage)
	_ = s.sut.GetTo(context.Background(), ch, s.id, uint32(end))
	s.events = s.drainChannel(ch)
	return nil
}

func (s *MemoryRepositorySuite) iRetrieveASubsetOfEventsFromTo(start, end int) error {
	ch := make(chan *evented.EventPage)
	_ = s.sut.GetFromTo(context.Background(), ch, s.id, uint32(start), uint32(end))
	s.events = s.drainChannel(ch)
	return nil
}

func (s *MemoryRepositorySuite) iRetrieveASubsetOfEventsStartingFromValue(start int) error {
	ch := make(chan *evented.EventPage)
	_ = s.sut.GetFrom(context.Background(), ch, s.id, uint32(start))
	s.events = s.drainChannel(ch)
	return nil
}

func (s *MemoryRepositorySuite) iRetrieveAllEvents() error {
	ch := make(chan *evented.EventPage)
	_ = s.sut.Get(context.Background(), ch, s.id)
	s.events = s.drainChannel(ch)
	return nil
}
