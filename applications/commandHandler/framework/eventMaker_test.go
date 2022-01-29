package framework

import (
	evented_proto "github.com/benjaminabbitt/evented/proto"
	"github.com/benjaminabbitt/evented/proto/gen/github.com/benjaminabbitt/evented/proto/evented/core"
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
	expected := &core.EventPage{
		Sequence:    &core.EventPage_Num{Num: 0},
		CreatedAt:   &timestamp.Timestamp{},
		Event:       nil,
		Synchronous: false,
	}
	o.Assert().Equal(expected, NewEventPage(0, false, nil))
}

func (o *EventMakerSuite) TestNewEmptyEventPage() {
	anyEmpty, _ := ptypes.MarshalAny(&empty.Empty{})
	page := &core.EventPage{
		Sequence:    &core.EventPage_Num{Num: 0},
		Synchronous: false,
		CreatedAt:   &timestamp.Timestamp{},
		Event:       anyEmpty,
	}
	o.Assert().Equal(page, NewEmptyEventPage(0, false))
}

func (o *EventMakerSuite) TestNewEventBook() {
	id, _ := uuid.NewRandom()
	protoUUID := evented_proto.UUIDToProto(id)
	pages := []*core.EventPage{NewEmptyEventPage(0, false)}
	eventBook := &core.EventBook{
		Cover: &core.Cover{
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
