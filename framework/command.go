package framework

import "github.com/golang/protobuf/ptypes/any"

type Command struct{
	Id string
	Sequence uint32
	Details *any.Any
}