package universal

import (
	"context"
	"fmt"
	eventedcore "github.com/benjaminabbitt/evented/proto/core"
	evented_projector "github.com/benjaminabbitt/evented/proto/projector"
	evented_query "github.com/benjaminabbitt/evented/proto/query"
	"github.com/benjaminabbitt/evented/repository/processed"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ProjectorCoordinator struct {
	Coordinator      *Coordinator
	Domain           string //Domain of the Source
	Log              *zap.SugaredLogger
	ProjectorClient  evented_projector.ProjectorClient
	EventQueryClient evented_query.EventQueryClient
	Processed        *processed.Processed
}

func (o *ProjectorCoordinator) HandleSync(ctx context.Context, eb *eventedcore.EventBook) (*eventedcore.Projection, error) {
	if eb.Cover.Domain != o.Domain {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("Event book Domain %s does not match projector configured Domain %s", eb.Cover.Domain, o.Domain))
	}
	o.Coordinator.RepairSequencing(ctx, eb, func(eb *eventedcore.EventBook) error {
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

func (o *ProjectorCoordinator) Handle(ctx context.Context, eb *eventedcore.EventBook) error {
	_, err := o.HandleSync(ctx, eb)
	return err
}
