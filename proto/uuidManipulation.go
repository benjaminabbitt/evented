package evented_proto

import (
	evented_core "github.com/benjaminabbitt/evented/proto/core"
	"github.com/google/uuid"
)

func ProtoToUUID(id evented_core.UUID) (u uuid.UUID, err error) {
	return uuid.FromBytes(id.Value)
}

func UUIDToProto(u uuid.UUID) (id evented_core.UUID) {
	return evented_core.UUID{
		Value: u[:],
	}
}
