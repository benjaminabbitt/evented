package repository

import "github.com/golang/protobuf/ptypes/any"

type StorageEvent struct{
	Id string
	Sequence uint32
	Details any.Any
}
