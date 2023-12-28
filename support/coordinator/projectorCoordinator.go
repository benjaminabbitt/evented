package coordinator

import (
	"context"
	"fmt"
	evented2 "github.com/benjaminabbitt/evented/generated/proto/github.com/benjaminabbitt/evented/proto/evented"
	"github.com/benjaminabbitt/evented/repository/processed"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ProjectorCoordinator struct {
	Coordinator      *Coordinator
	Domain           string //Domain of the Source
	ProjectorClient  evented2.ProjectorClient
	EventQueryClient evented2.EventQueryClient
	Processed        *processed.Processed
	Log              *zap.SugaredLogger
}

func (o ProjectorCoordinator) HandleSync(ctx context.Context, eb *evented2.EventBook) (*evented2.Projection, error) {
	if eb.Cover.Domain != o.Domain {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("Event book Domain %s does not match sample-projector configured Domain %s", eb.Cover.Domain, o.Domain))
	}
	o.Coordinator.RepairSequencing(ctx, eb, func(eb *evented2.EventBook) error {
		_, err := o.ProjectorClient.Handle(ctx, eb)
		return err
	})

	reb, err := o.ProjectorClient.HandleSync(ctx, eb)
	if err != nil {
		o.Log.Error(err)
	}
	o.Coordinator.MarkProcessed(ctx, eb)
	return reb, err
}

func (o ProjectorCoordinator) Handle(ctx context.Context, eb *evented2.EventBook) error {
	_, err := o.HandleSync(ctx, eb)
	return err
}
