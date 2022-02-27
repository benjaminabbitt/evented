package saga

import (
	evented_proto "github.com/benjaminabbitt/evented/proto"
	"github.com/benjaminabbitt/evented/proto/gen/github.com/benjaminabbitt/evented/proto/evented"

	"github.com/benjaminabbitt/evented/support"
	"github.com/benjaminabbitt/evented/support/grpcWithInterceptors"
	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"golang.org/x/net/context"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func NewPlaceholderSagaLogic(log *zap.SugaredLogger, tracer *opentracing.Tracer) PlaceholderSagaLogic {
	return PlaceholderSagaLogic{
		log:    log,
		tracer: tracer,
	}
}

type PlaceholderSagaLogic struct {
	evented.UnimplementedSagaServer
	eventDomain string
	log         *zap.SugaredLogger
	tracer      *opentracing.Tracer
}

func (o *PlaceholderSagaLogic) Handle(ctx context.Context, in *evented.EventBook) (*evented.EventBook, error) {
	return o.HandleSync(ctx, in)
}

func (o *PlaceholderSagaLogic) HandleSync(ctx context.Context, in *evented.EventBook) (*evented.EventBook, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		o.log.Error(err)
	}
	root := evented_proto.UUIDToProto(id)
	cover := evented.Cover{
		Domain: o.eventDomain,
		Root:   &root,
	}
	eb := &evented.EventBook{
		Cover: &cover,
		Pages: []*evented.EventPage{&evented.EventPage{
			Sequence:    &evented.EventPage_Force{Force: true},
			CreatedAt:   &timestamppb.Timestamp{},
			Event:       nil,
			Synchronous: false,
		}},
		Snapshot: nil,
	}
	return eb, nil
}

func (o *PlaceholderSagaLogic) Listen(port uint) {
	lis := support.CreateListener(port, o.log)
	grpcServer := grpcWithInterceptors.GenerateConfiguredServer(o.log.Desugar(), *o.tracer)

	evented.RegisterSagaServer(grpcServer, o)
	err := grpcServer.Serve(lis)
	if err != nil {
		o.log.Error(err)
	}
}