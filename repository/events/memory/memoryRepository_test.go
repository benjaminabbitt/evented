package memory

import (
	"context"
	evented_core "github.com/benjaminabbitt/evented/proto/evented/core"
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

var log *zap.SugaredLogger
var sut EventRepoMemory
var id uuid.UUID
var events []*evented_core.EventPage

func InitializeTestSuite(ctx *godog.TestSuiteContext) {
	log = support.Log()
}

func InitializeScenario(s *godog.ScenarioContext) {
	sut, _ = NewEventRepoMemory(log)
	s.Step(`^I should be able to retrieve it by its coordinates:$`, iShouldBeAbleToRetrieveItByItsCoordinates)
	s.Step(`^I store the event:$`, iStoreTheEvent)
	s.Step(`^a populated database:$`, aPopulatedDatabase)
	s.Step(`^I should get these events:$`, iShouldGetTheseEvents)
	s.Step(`^I retrieve a subset of events ending at event (\d+)$`, iRetrieveASubsetOfEventsEndingAtEvent)
	s.Step(`^I retrieve a subset of events from (\d+) to (\d+)$`, iRetrieveASubsetOfEventsFromTo)
	s.Step(`^I retrieve a subset of events starting from value (\d+)$`, iRetrieveASubsetOfEventsStartingFromValue)
	s.Step(`^I retrieve all events$`, iRetrieveAllEvents)
}

func iShouldBeAbleToRetrieveItByItsCoordinates(arg1 *messages.PickleStepArgument_PickleTable) error {
	id, events := extractPickleTableToEvents(arg1)
	ch := make(chan *evented_core.EventPage)
	_ = sut.Get(context.Background(), ch, id)
	return cucumber.AssertExpectedAndActual(assert.Equal, events[0], <-ch, "", "")
}

func extractPickleTableToEvents(arg *messages.PickleStepArgument_PickleTable) (id uuid.UUID, events []*evented_core.EventPage) {
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

		event := &evented_core.EventPage{
			CreatedAt:   ts,
			Event:       nil,
			Synchronous: false,
		}
		if force {
			event.Sequence = &evented_core.EventPage_Force{}
		} else {
			event.Sequence = &evented_core.EventPage_Num{Num: sequence}
		}
		events = append(events, event)
	}
	return id, events
}

func iStoreTheEvent(arg1 *messages.PickleStepArgument_PickleTable) error {
	id, events = extractPickleTableToEvents(arg1)
	_ = sut.Add(context.Background(), id, events)
	return nil
}

func aPopulatedDatabase(arg1 *messages.PickleStepArgument_PickleTable) error {
	id, events = extractPickleTableToEvents(arg1)
	_ = sut.Add(context.Background(), id, events)
	return nil
}

func iShouldGetTheseEvents(arg1 *messages.PickleStepArgument_PickleTable) error {
	_, expectedEvents := extractPickleTableToEvents(arg1)
	return cucumber.AssertExpectedAndActual(assert.Equal, expectedEvents, events, "", "")
}

func drainChannel(ch chan *evented_core.EventPage) (pages []*evented_core.EventPage) {
	for page := range ch {
		pages = append(pages, page)
	}
	return pages
}

func iRetrieveASubsetOfEventsEndingAtEvent(end int) error {
	ch := make(chan *evented_core.EventPage)
	_ = sut.GetTo(context.Background(), ch, id, uint32(end))
	events = drainChannel(ch)
	return nil
}

func iRetrieveASubsetOfEventsFromTo(start, end int) error {
	ch := make(chan *evented_core.EventPage)
	_ = sut.GetFromTo(context.Background(), ch, id, uint32(start), uint32(end))
	events = drainChannel(ch)
	return nil
}

func iRetrieveASubsetOfEventsStartingFromValue(start int) error {
	ch := make(chan *evented_core.EventPage)
	_ = sut.GetFrom(context.Background(), ch, id, uint32(start))
	events = drainChannel(ch)
	return nil
}

func iRetrieveAllEvents() error {
	ch := make(chan *evented_core.EventPage)
	_ = sut.Get(context.Background(), ch, id)
	events = drainChannel(ch)
	return nil
}
