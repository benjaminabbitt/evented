package framework

import (
	"github.com/benjaminabbitt/evented/generated/proto/github.com/benjaminabbitt/evented/proto/evented"
	evented_proto "github.com/benjaminabbitt/evented/proto"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func NewEventPage(sequence uint32, sync bool, eventDetails *anypb.Any) *evented.EventPage {
	return &evented.EventPage{
		Sequence:    &evented.EventPage_Num{Num: sequence},
		Synchronous: sync,
		CreatedAt:   &timestamppb.Timestamp{},
		Event:       eventDetails,
	}
}

func NewEmptyEventPage(sequence uint32, sync bool) *evented.EventPage {
	anyEmpty, _ := anypb.New(&emptypb.Empty{})
	return NewEventPage(sequence, sync, anyEmpty)
}

func NewEventBook(id uuid.UUID, domain string, events []*evented.EventPage, snapshot *evented.Snapshot) *evented.EventBook {
	protoUUID := evented_proto.UUIDToProto(id)
	return &evented.EventBook{
		Cover: &evented.Cover{
			Domain: domain,
			Root:   &protoUUID,
		},
		Pages:    events,
		Snapshot: snapshot,
	}
}
