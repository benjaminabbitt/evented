package evented_proto

import (
	"github.com/benjaminabbitt/evented/proto/gen/github.com/benjaminabbitt/evented/proto/evented"
	"github.com/google/uuid"
)

func ProtoToUUID(id *evented.UUID) (u uuid.UUID, err error) {
	return uuid.FromBytes(id.Value)
}

func UUIDToProto(u uuid.UUID) (id evented.UUID) {
	return evented.UUID{
		Value: u[:],
	}
}
