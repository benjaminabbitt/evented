package evented_proto

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"testing"
)

type UUIDManipSuite struct {
	suite.Suite
}

func (s *UUIDManipSuite) Test_Protobuf_Ser_Deser() {
	input, _ := uuid.NewRandom()
	proto := UUIDToProto(input)
	output, _ := ProtoToUUID(&proto)
	s.EqualValues(input.String(), output.String())
}

func TestUUIDManipSuite(t *testing.T) {
	suite.Run(t, new(UUIDManipSuite))
}
