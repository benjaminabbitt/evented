package universal

import (
	"fmt"
	"github.com/benjaminabbitt/evented/proto/gen/github.com/benjaminabbitt/evented/proto/evented/business"

	"github.com/benjaminabbitt/evented/proto/gen/github.com/benjaminabbitt/evented/proto/evented/saga"
	"go.uber.org/zap"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type SagaCoordinator struct {
	Coordinator         *Coordinator
	Domain              string
	SagaClient          saga.SagaClient
	OtherCommandHandler business.BusinessCoordinatorClient
	Log                 *zap.SugaredLogger
}

func (o *SagaCoordinator) HandleSync(ctx context.Context, eb *evented.EventBook) (*evented.SynchronousProcessingResponse, error) {
	if eb.Cover.Domain != o.Domain {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("Event book Domain %s does not match saga configured Domain %s", eb.Cover.Domain, o.Domain))
	}

	o.Coordinator.RepairSequencing(ctx, eb, func(eb *evented.EventBook) error {
		_, err := o.SagaClient.Handle(ctx, eb)
		return err
	})

	sagaResponseBooks, err := o.SagaClient.HandleSync(ctx, eb)
	if err != nil {
		o.Log.Error(err)
	}
	o.Coordinator.MarkProcessed(ctx, eb)
	commandHandlerResponse, err := o.OtherCommandHandler.Record(ctx, sagaResponseBooks)
	if err != nil {
		o.Log.Error(err)
	}
	commandHandlerResponse.Books = append(commandHandlerResponse.Books, sagaResponseBooks)
	return commandHandlerResponse, err
}

func (o *SagaCoordinator) Handle(ctx context.Context, eb *evented.EventBook) (err error) {
	_, err = o.HandleSync(ctx, eb)
	return err
}
