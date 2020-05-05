package eventQueryServer

import (
	"context"
	evented_core "github.com/benjaminabbitt/evented/proto/evented/core"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/metadata"
)

type MockGetEventsServer struct {
	mock.Mock
}

func (o MockGetEventsServer) Send(book *evented_core.EventBook) error {
	args := o.Called(book)
	return args.Error(0)
}

func (o MockGetEventsServer) SetHeader(md metadata.MD) error {
	args := o.Called(md)
	return args.Error(0)
}

func (o MockGetEventsServer) SendHeader(md metadata.MD) error {
	args := o.Called(md)
	return args.Error(0)
}

func (o MockGetEventsServer) SetTrailer(md metadata.MD) {
	o.Called(md)
}

func (o MockGetEventsServer) Context() context.Context {
	args := o.Called()
	return args.Get(0).(context.Context)
}

func (o MockGetEventsServer) SendMsg(m interface{}) error {
	args := o.Called(m)
	return args.Error(0)
}

func (o MockGetEventsServer) RecvMsg(m interface{}) error {
	args := o.Called(m)
	return args.Error(0)
}
