package mongo

import "go.mongodb.org/mongo-driver/bson"

type SingleResult interface {
	Decode(v interface{}) error
	DecodeBytes() (bson.Raw, error)
	Err() error
}
