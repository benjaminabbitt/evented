package coordinator

import (
	"fmt"
	evented2 "github.com/benjaminabbitt/evented/generated/proto/github.com/benjaminabbitt/evented/proto/evented"
	"go.uber.org/zap"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type SagaCoordinator struct {
	Coordinator         *Coordinator
	Domain              string
	SagaClient          evented2.SagaClient
	OtherCommandHandler evented2.BusinessCoordinatorClient
	Log                 *zap.SugaredLogger
}

func (o *SagaCoordinator) HandleSync(ctx context.Context, eb *evented2.EventBook) (*evented2.SynchronousProcessingResponse, error) {
	if eb.Cover.Domain != o.Domain {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("Event book Domain %s does not match sample-saga configured Domain %s", eb.Cover.Domain, o.Domain))
	}

	o.Coordinator.RepairSequencing(ctx, eb, func(eb *evented2.EventBook) error {
		_, err := o.SagaClient.Handle(ctx, eb)
		return err
	})

	sagaResponseBooks, err := o.SagaClient.HandleSync(ctx, eb)
	if err != nil {
		o.Log.Error(err)
	}
	o.Coordinator.MarkProcessed(ctx, eb)

	commandHandlerResponse := &evented2.SynchronousProcessingResponse{
		Books:       []*evented2.EventBook{},
		Projections: []*evented2.Projection{},
	}

	for _, book := range sagaResponseBooks.Books {
		otherCommandHandlerResponse, err := o.OtherCommandHandler.Record(ctx, book)
		if err != nil {
			o.Log.Error(err)
		}
		commandHandlerResponse.Books = append(commandHandlerResponse.Books, otherCommandHandlerResponse.Books...)
	}
	return commandHandlerResponse, err
}

func (o *SagaCoordinator) Handle(ctx context.Context, eb *evented2.EventBook) (err error) {
	_, err = o.HandleSync(ctx, eb)
	return err
}
