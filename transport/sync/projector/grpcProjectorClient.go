package projector

import (
	"context"
	"github.com/benjaminabbitt/evented/proto/core"
	evented_projector "github.com/benjaminabbitt/evented/proto/projector"
	"go.uber.org/zap"
)

type GrpcProjector struct {
	log    *zap.SugaredLogger
	client evented_projector.ProjectorClient
}

func (o GrpcProjector) HandleSync(ctx context.Context, evts *evented_core.EventBook) (projection *evented_core.Projection, err error) {
	return o.client.HandleSync(ctx, evts)
}
