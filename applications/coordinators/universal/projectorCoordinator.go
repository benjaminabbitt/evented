package universal

import (
	"context"
	"fmt"
	eventedcore "github.com/benjaminabbitt/evented/proto/evented/core"
	eventedprojector "github.com/benjaminabbitt/evented/proto/evented/projector"
	eventedquery "github.com/benjaminabbitt/evented/proto/evented/query"
	"github.com/benjaminabbitt/evented/repository/processed"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ProjectorCoordinator struct {
	Coordinator      *Coordinator
	Domain           string //Domain of the Source
	ProjectorClient  eventedprojector.ProjectorClient
	EventQueryClient eventedquery.EventQueryClient
	Processed        *processed.Processed
	Log              *zap.SugaredLogger
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
