package projector

import (
	evented2 "github.com/benjaminabbitt/evented/generated/proto/github.com/benjaminabbitt/evented/proto/evented"
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
	evented2.UnimplementedProjectorServer
	eventDomain string
	log         *zap.SugaredLogger
	tracer      *opentracing.Tracer
}

func (o PlaceholderProjectorLogic) Handle(ctx context.Context, in *evented2.EventBook) (*empty.Empty, error) {
	o.log.Infow("In Handle", in.String())
	_, err := o.HandleSync(ctx, in)
	return &empty.Empty{}, err
}

func (o PlaceholderProjectorLogic) HandleSync(ctx context.Context, in *evented2.EventBook) (projection *evented2.Projection, err error) {
	o.log.Infow("In HandleSync", in.String())
	lastSequenceIndex := len(in.Pages) - 1
	lastSequence := in.Pages[lastSequenceIndex].Sequence.(*evented2.EventPage_Num).Num
	projection = &evented2.Projection{
		Cover:      in.Cover,
		Projector:  "test",
		Sequence:   lastSequence,
		Projection: nil,
	}
	o.log.Infow("Returning projection...")
	return projection, nil
}

func (o PlaceholderProjectorLogic) Listen(port uint) {
	lis := support.CreateListener(port, o.log)
	grpcServer := grpcWithInterceptors.GenerateConfiguredServer(o.log.Desugar(), *o.tracer)

	evented2.RegisterProjectorServer(grpcServer, o)
	err := grpcServer.Serve(lis)
	if err != nil {
		o.log.Error(err)
	}
}
