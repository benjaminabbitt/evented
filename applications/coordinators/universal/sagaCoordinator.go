package universal

import (
	"fmt"
	eventedcore "github.com/benjaminabbitt/evented/proto/evented/core"
	eventedsaga "github.com/benjaminabbitt/evented/proto/evented/saga"
	"go.uber.org/zap"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type SagaCoordinator struct {
	Coordinator         *Coordinator
	Domain              string
	SagaClient          eventedsaga.SagaClient
	OtherCommandHandler eventedcore.CommandHandlerClient
	Log                 *zap.SugaredLogger
}

func (o *SagaCoordinator) HandleSync(ctx context.Context, eb *eventedcore.EventBook) (*eventedcore.SynchronousProcessingResponse, error) {
	if eb.Cover.Domain != o.Domain {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("Event book Domain %s does not match saga configured Domain %s", eb.Cover.Domain, o.Domain))
	}

	o.Coordinator.RepairSequencing(ctx, eb, func(eb *eventedcore.EventBook) error {
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

func (o *SagaCoordinator) Handle(ctx context.Context, eb *eventedcore.EventBook) (err error) {
	_, err = o.HandleSync(ctx, eb)
	return err
}
