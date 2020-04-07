package framework

import (
	evented_proto "github.com/benjaminabbitt/evented/proto"
	evented_core "github.com/benjaminabbitt/evented/proto/core"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/google/uuid"
)

func NewEventPage(sequence uint32, sync bool, eventDetails any.Any) *evented_core.EventPage {
	return &evented_core.EventPage{
		Sequence:    &evented_core.EventPage_Num{Num:sequence},
		Synchronous: sync,
		CreatedAt:   &timestamp.Timestamp{},
		Event:       &eventDetails,
	}
}

func NewEventBook(id uuid.UUID, domain string, events []*evented_core.EventPage, snapshot *evented_core.Snapshot) *evented_core.EventBook {
	protoUUID := evented_proto.UUIDToProto(id)
	return &evented_core.EventBook{
		Cover: &evented_core.Cover{
			Domain: domain,
			Root:   &protoUUID,
		},
		Pages:    events,
		Snapshot: snapshot,
	}
}
