package framework

import (
	evented_proto "github.com/benjaminabbitt/evented/proto"
	evented_core "github.com/benjaminabbitt/evented/proto/core"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"testing"
)

type EventMakerSuite struct {
	suite.Suite
}

func (o *EventMakerSuite) TestNewEventPage() {
	expected := &evented_core.EventPage{
		Sequence:    &evented_core.EventPage_Num{Num: 0},
		CreatedAt:   &timestamp.Timestamp{},
		Event:       nil,
		Synchronous: false,
	}
	o.Assert().Equal(expected, NewEventPage(0, false, nil))
}

func (o *EventMakerSuite) TestNewEmptyEventPage() {
	anyEmpty, _ := ptypes.MarshalAny(&empty.Empty{})
	page := &evented_core.EventPage{
		Sequence:    &evented_core.EventPage_Num{Num: 0},
		Synchronous: false,
		CreatedAt:   &timestamp.Timestamp{},
		Event:       anyEmpty,
	}
	o.Assert().Equal(page, NewEmptyEventPage(0, false))
}

func (o *EventMakerSuite) TestNewEventBook() {
	id, _ := uuid.NewRandom()
	protoUUID := evented_proto.UUIDToProto(id)
	pages := []*evented_core.EventPage{NewEmptyEventPage(0, false)}
	eventBook := &evented_core.EventBook{
		Cover: &evented_core.Cover{
			Domain: "",
			Root:   &protoUUID,
		},
		Pages:    pages,
		Snapshot: nil,
	}
	o.Assert().Equal(eventBook, NewEventBook(id, "", pages, nil))
}
func TestEventMakerSuite(t *testing.T) {
	suite.Run(t, new(EventMakerSuite))
}
