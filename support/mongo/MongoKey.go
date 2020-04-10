package mongosupport

import "github.com/google/uuid"

func RootToMongo(id uuid.UUID) (idBytes [12]byte, err error) {
	idByteSlice, err := id.MarshalBinary()
	if err != nil {
		return idBytes, err
	}
	copy(idBytes[:], idByteSlice[:12])
	return idBytes, nil
}
