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
	"strconv"
	"time"
)

func InitializeTestSuite(ctx *godog.TestSuiteContext) {
	log := support.Log()
	sut, _ = NewEventRepoMemory(log)
}

var sut EventRepoMemory

func InitializeScenario(scenarioContext *godog.ScenarioContext) {
	scenarioContext.Step(`^I should be able to retrieve it by its coordinates:$`, iShouldBeAbleToRetrieveItByItsCoordinates)
	scenarioContext.Step(`^I store the event:$`, iStoreTheEvent)
}

func iShouldBeAbleToRetrieveItByItsCoordinates(arg1 *messages.PickleStepArgument_PickleTable) error {
	id, events := extractPickleTableToEvents(arg1)
	ch := make(chan *evented_core.EventPage)
	_ = sut.Get(context.Background(), ch, id)
	return cucumber.AssertExpectedAndActual(assert.Equal, events[0], <-ch, "", "")
}

func extractPickleTableToEvents(arg *messages.PickleStepArgument_PickleTable) (id uuid.UUID, events []*evented_core.EventPage) {
	for _, row := range arg.GetRows() {
		var sequence uint32
		var force bool
		var ts *timestamppb.Timestamp
		for i, cell := range row.GetCells() {
			switch i {
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
	id, events := extractPickleTableToEvents(arg1)
	_ = sut.Add(context.Background(), id, events)
	return nil
}
