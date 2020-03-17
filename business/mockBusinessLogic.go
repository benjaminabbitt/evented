package business

import (
	eventedcore "github.com/benjaminabbitt/evented/proto/core"
	"github.com/stretchr/testify/mock"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type MockBusinessLogicClient struct {
	mock.Mock
}

func (c *MockBusinessLogicClient) Handle(ctx context.Context, in *eventedcore.ContextualCommand, opts ...grpc.CallOption) (*eventedcore.EventBook, error){
	c.Called(ctx, in, opts)
	return &eventedcore.EventBook{
		Cover:    in.Command.Cover,
		Pages:    nil,
		Snapshot: nil,
	}, nil
}
