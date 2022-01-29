package projector

import (
	"github.com/benjaminabbitt/evented/proto/gen/github.com/benjaminabbitt/evented/proto/evented/core"
	"github.com/benjaminabbitt/evented/proto/gen/github.com/benjaminabbitt/evented/proto/evented/projector"
	"github.com/benjaminabbitt/evented/support"
	"github.com/benjaminabbitt/evented/support/grpcWithInterceptors"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"golang.org/x/net/context"
)

func NewPlaceholderProjectorLogic(log *zap.SugaredLogger, tracer *opentracing.Tracer) PlaceholderProjectorLogic {
	return PlaceholderProjectorLogic{
		log:    log,
		tracer: tracer,
	}
}

type PlaceholderProjectorLogic struct {
	projector.UnimplementedProjectorServer
	eventDomain string
	log         *zap.SugaredLogger
	tracer      *opentracing.Tracer
}

func (o *PlaceholderProjectorLogic) Handle(ctx context.Context, in *core.EventBook) (*empty.Empty, error) {
	_, err := o.HandleSync(ctx, in)
	return &empty.Empty{}, err
}

func (o *PlaceholderProjectorLogic) HandleSync(ctx context.Context, in *core.EventBook) (*core.Projection, error) {
	lastSequenceIndex := len(in.Pages) - 1
	lastSequence := in.Pages[lastSequenceIndex].Sequence.(*core.EventPage_Num).Num
	projection := &core.Projection{
		Cover:      in.Cover,
		Projector:  "test",
		Sequence:   lastSequence,
		Projection: nil,
	}
	return projection, nil
}

func (o *PlaceholderProjectorLogic) Listen(port uint) {
	lis := support.CreateListener(port, o.log)
	grpcServer := grpcWithInterceptors.GenerateConfiguredServer(o.log.Desugar(), *o.tracer)

	projector.RegisterProjectorServer(grpcServer, o)
	err := grpcServer.Serve(lis)
	if err != nil {
		o.log.Error(err)
	}
}
