package projector

import (
	eventedcore "github.com/benjaminabbitt/evented/proto/core"
	evented_projector "github.com/benjaminabbitt/evented/proto/projector"
	"github.com/benjaminabbitt/evented/support"
	"github.com/benjaminabbitt/evented/support/grpcWithInterceptors"
	"github.com/golang/protobuf/ptypes/empty"
	"go.uber.org/zap"
	"golang.org/x/net/context"
)

func NewPlaceholderProjectorLogic(log *zap.SugaredLogger) PlaceholderProjectorLogic {
	return PlaceholderProjectorLogic{
		log: log,
	}
}

type PlaceholderProjectorLogic struct {
	evented_projector.UnimplementedProjectorServer
	eventDomain string
	log         *zap.SugaredLogger
}

func (o *PlaceholderProjectorLogic) Handle(ctx context.Context, in *eventedcore.EventBook) (*empty.Empty, error) {
	_, err := o.HandleSync(ctx, in)
	return &empty.Empty{}, err
}

func (o *PlaceholderProjectorLogic) HandleSync(ctx context.Context, in *eventedcore.EventBook) (*eventedcore.Projection, error) {
	lastSequenceIndex := len(in.Pages) - 1
	lastSequence := in.Pages[lastSequenceIndex].Sequence.(*eventedcore.EventPage_Num).Num
	projection := &eventedcore.Projection{
		Cover:      in.Cover,
		Projector:  "test",
		Sequence:   lastSequence,
		Projection: nil,
	}
	return projection, nil
}

func (o *PlaceholderProjectorLogic) Listen(port uint) {
	lis := support.CreateListener(port, o.log)
	grpcServer := grpcWithInterceptors.GenerateConfiguredServer(o.log.Desugar())

	evented_projector.RegisterProjectorServer(grpcServer, o)
	err := grpcServer.Serve(lis)
	if err != nil {
		o.log.Error(err)
	}
}