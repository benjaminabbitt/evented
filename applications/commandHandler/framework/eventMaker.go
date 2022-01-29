package framework

import (
	evented_proto "github.com/benjaminabbitt/evented/proto"
	"github.com/benjaminabbitt/evented/proto/gen/github.com/benjaminabbitt/evented/proto/evented/core"
	"github.com/google/uuid"
	anypb "google.golang.org/protobuf/types/known/anypb"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

func NewEventPage(sequence uint32, sync bool, eventDetails *anypb.Any) *core.EventPage {
	return &core.EventPage{
		Sequence:    &core.EventPage_Num{Num: sequence},
		Synchronous: sync,
		CreatedAt:   &timestamppb.Timestamp{},
		Event:       eventDetails,
	}
}

func NewEmptyEventPage(sequence uint32, sync bool) *core.EventPage {
	anyEmpty, _ := anypb.New(&emptypb.Empty{})
	return &core.EventPage{
		Sequence:    &core.EventPage_Num{Num: sequence},
		Synchronous: sync,
		CreatedAt:   &timestamppb.Timestamp{},
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
