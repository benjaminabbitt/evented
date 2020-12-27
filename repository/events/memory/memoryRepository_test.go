package memory

import (
	"context"
	evented_core "github.com/benjaminabbitt/evented/proto/evented/core"
	"github.com/benjaminabbitt/evented/support"
	"github.com/benjaminabbitt/evented/support/cucumber"
	"github.com/cucumber/godog"
	"github.com/cucumber/messages-go/v10"
	"github.com/golang/protobuf/ptypes"
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
	scenarioContext.Step(`^I store the sample event:$`, iStoreTheSampleEvent)
	scenarioContext.Step(`^that we\'re working in the coordinates of a domain$`, thatWereWorkingInTheCoordinatesOfADomain)
}

func iShouldBeAbleToRetrieveItByItsCoordinates(arg1 *messages.PickleStepArgument_PickleTable) error {
	for _, row := range arg1.GetRows() {
		var id uuid.UUID
		var sequence uint32
		for i, cell := range row.GetCells() {
			switch i {
			case 0:
				id, _ = uuid.Parse(cell.Value)
			case 1:
				sequence, _ = cucumber.Uint64ToUint32WithErrorPassthrough(strconv.ParseUint(cell.Value, 10, 32))
			}
		}
		ch := make(chan *evented_core.EventPage)
		_ = sut.Get(context.Background(), ch, id)
		return cucumber.AssertExpectedAndActual(assert.Equal, &evented_core.EventPage_Num{Num: sequence}, (<-ch).Sequence, "", "")
	}
	return nil
}

func iStoreTheSampleEvent(arg1 *messages.PickleStepArgument_PickleTable) error {
	for _, row := range arg1.GetRows() {
		var id uuid.UUID
		var sequence uint32
		for i, cell := range row.GetCells() {
			switch i {
			case 0:
				id, _ = uuid.Parse(cell.Value)
			case 1:
				sequence, _ = cucumber.Uint64ToUint32WithErrorPassthrough(strconv.ParseUint(cell.Value, 10, 32))
			}
		}
		ts, _ := ptypes.TimestampProto(time.Now())

		page := &evented_core.EventPage{
			Sequence:    &evented_core.EventPage_Num{Num: sequence},
			CreatedAt:   ts,
			Event:       nil,
			Synchronous: false,
		}

		_ = sut.Add(context.Background(), id, []*evented_core.EventPage{page})
	}
	return nil
}

func thatWereWorkingInTheCoordinatesOfADomain() error {
	return nil
}
