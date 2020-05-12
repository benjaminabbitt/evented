package client

import (
	"context"
	eventedcore "github.com/benjaminabbitt/evented/proto/evented/core"
	"github.com/stretchr/testify/mock"
)

type MockClient struct {
	mock.Mock
}

func (client *MockClient) Handle(ctx context.Context, command *eventedcore.ContextualCommand) (events *eventedcore.EventBook, err error) {
	args := client.Called(ctx, command)
	return args.Get(0).(*eventedcore.EventBook), args.Error(1)
}
