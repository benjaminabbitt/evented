package saga

import (
	evented_proto "github.com/benjaminabbitt/evented/proto"
	eventedcore "github.com/benjaminabbitt/evented/proto/evented/core"
	"github.com/benjaminabbitt/evented/proto/evented/saga/saga"
	"github.com/benjaminabbitt/evented/support"
	"github.com/benjaminabbitt/evented/support/grpcWithInterceptors"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"golang.org/x/net/context"
)

func NewPlaceholderSagaLogic(log *zap.SugaredLogger, tracer *opentracing.Tracer) PlaceholderSagaLogic {
	return PlaceholderSagaLogic{
		log:    log,
		tracer: tracer,
	}
}

type PlaceholderSagaLogic struct {
	saga.UnimplementedSagaServer
	eventDomain string
	log         *zap.SugaredLogger
	tracer      *opentracing.Tracer
}

func (o *PlaceholderSagaLogic) Handle(ctx context.Context, in *eventedcore.EventBook) (*eventedcore.EventBook, error) {
	return o.HandleSync(ctx, in)
}

func (o *PlaceholderSagaLogic) HandleSync(ctx context.Context, in *eventedcore.EventBook) (*eventedcore.EventBook, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		o.log.Error(err)
	}
	root := evented_proto.UUIDToProto(id)
	cover := eventedcore.Cover{
		Domain: o.eventDomain,
		Root:   &root,
	}
	eb := &eventedcore.EventBook{
		Cover: &cover,
		Pages: []*eventedcore.EventPage{&eventedcore.EventPage{
			Sequence:    &eventedcore.EventPage_Force{Force: true},
			CreatedAt:   &timestamp.Timestamp{},
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

	saga.RegisterSagaServer(grpcServer, o)
	err := grpcServer.Serve(lis)
	if err != nil {
		o.log.Error(err)
	}
}
