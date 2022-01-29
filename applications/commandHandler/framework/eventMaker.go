package framework

import (
	evented_proto "github.com/benjaminabbitt/evented/proto"
	"github.com/benjaminabbitt/evented/proto/gen/github.com/benjaminabbitt/evented/proto/evented/core"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/google/uuid"
)

func NewEventPage(sequence uint32, sync bool, eventDetails *any.Any) *core.EventPage {
	return &core.EventPage{
		Sequence:    &core.EventPage_Num{Num: sequence},
		Synchronous: sync,
		CreatedAt:   &timestamp.Timestamp{},
		Event:       eventDetails,
	}
}

func NewEmptyEventPage(sequence uint32, sync bool) *core.EventPage {
	anyEmpty, _ := ptypes.MarshalAny(&empty.Empty{})
	return &core.EventPage{
		Sequence:    &core.EventPage_Num{Num: sequence},
		Synchronous: sync,
		CreatedAt:   &timestamp.Timestamp{},
		Event:       anyEmpty,
	}
}

func NewEventBook(id uuid.UUID, domain string, events []*core.EventPage, snapshot *core.Snapshot) *core.EventBook {
	protoUUID := evented_proto.UUIDToProto(id)
	return &core.EventBook{
		Cover: &core.Cover{
			Domain: domain,
			Root:   &protoUUID,
		},
		Pages:    events,
		Snapshot: snapshot,
	}
}
