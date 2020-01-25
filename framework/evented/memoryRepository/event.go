package memoryRepository

import "github.com/golang/protobuf/ptypes/any"

type Event struct{
	id string
	sequence uint32
	details any.Any
}
