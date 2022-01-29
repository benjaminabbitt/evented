package projector

import (
	"github.com/benjaminabbitt/evented/proto/gen/github.com/benjaminabbitt/evented/proto/evented"
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
	evented.UnimplementedProjectorServer
	eventDomain string
	log         *zap.SugaredLogger
	tracer      *opentracing.Tracer
}

func (o *PlaceholderProjectorLogic) Handle(ctx context.Context, in *evented.EventBook) (*empty.Empty, error) {
	_, err := o.HandleSync(ctx, in)
	return &empty.Empty{}, err
}

func (o *PlaceholderProjectorLogic) HandleSync(ctx context.Context, in *evented.EventBook) (*evented.Projection, error) {
	lastSequenceIndex := len(in.Pages) - 1
	lastSequence := in.Pages[lastSequenceIndex].Sequence.(*evented.EventPage_Num).Num
	projection := &evented.Projection{
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

	evented.RegisterProjectorServer(grpcServer, o)
	err := grpcServer.Serve(lis)
	if err != nil {
		o.log.Error(err)
	}
}
