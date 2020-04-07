package businessLogic

import (
	"context"
	"github.com/benjaminabbitt/evented"
	"github.com/benjaminabbitt/evented/proto/core"
	evented_query "github.com/benjaminabbitt/evented/proto/query"
	"github.com/benjaminabbitt/evented/repository/events"
	"github.com/benjaminabbitt/evented/support"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func NewEventQueryServer(repos events.EventRepository, log *zap.SugaredLogger, errh *evented.ErrLogger) EventQueryServer {
	return EventQueryServer{
		repos: repos,
		log:  log,
		errh: errh,
	}
}

type EventQueryServer struct {
	evented_query.UnimplementedEventQueryServer
	repos events.EventRepository
	log  *zap.SugaredLogger
	errh *evented.ErrLogger
}

func (server *EventQueryServer) GetEventBook(ctx context.Context, req *evented_query.Query) (*evented_core.EventBook, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetEventBook not implemented")
}
func (server *EventQueryServer) GetNextSequence(ctx context.Context, req *evented_core.UUID) (*evented_query.NextSequence, error) {
	id, err := uuid.ParseBytes(req.Value)
	if err != nil {
		return nil, err
	}
	seq, err := server.repos.GetNextSequence(id)
	if err != nil {
		return nil, err
	}
	return seq, nil
}

func (server *EventQueryServer) Listen(port uint16) {
	lis := support.CreateListener(port, server.errh)
	grpcServer := grpc.NewServer()

	evented_query.RegisterEventQueryServer(grpcServer, server)
	err := grpcServer.Serve(lis)
	server.errh.LogIfErr(err, "Failed starting server")
}
