package mongosupport

import (
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson"
)

type MockSingleResult struct {
	mock.Mock
}

func (o *MockSingleResult) Decode(v interface{}) error {
	args := o.Called(v)
	return args.Error(0)
}

func (o *MockSingleResult) DecodeBytes() (bson.Raw, error) {
	args := o.Called()
	return args.Get(0).(bson.Raw), args.Error(1)
}

func (o *MockSingleResult) Err() error {
	args := o.Called()
	return args.Error(0)
}
