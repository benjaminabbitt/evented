package transport

import "github.com/golang/protobuf/ptypes/any"

type TransportEvent struct{
	Id string
	Sequence uint32
	Details any.Any
}
