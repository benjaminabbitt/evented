package evented_proto

import (
	core "github.com/benjaminabbitt/evented/proto/evented/core"
	"github.com/google/uuid"
)

func ProtoToUUID(id *core.UUID) (u uuid.UUID, err error) {
	return uuid.FromBytes(id.Value)
}

func UUIDToProto(u uuid.UUID) (id core.UUID) {
	return core.UUID{
		Value: u[:],
	}
}
