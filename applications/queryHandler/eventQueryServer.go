package main

import (
	"context"
	"github.com/benjaminabbitt/evented"
	evented_proto "github.com/benjaminabbitt/evented/proto"
	"github.com/benjaminabbitt/evented/proto/core"
	evented_query "github.com/benjaminabbitt/evented/proto/query"
	"github.com/benjaminabbitt/evented/repository/eventBook"
	"github.com/benjaminabbitt/evented/support"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func NewEventQueryServer(repos eventBook.Repository, log *zap.SugaredLogger, errh *evented.ErrLogger) EventQueryServer {
	return EventQueryServer{
		repos: repos,
		log:   log,
		errh:  errh,
	}
}

type EventQueryServer struct {
	evented_query.UnimplementedEventQueryServer
	repos eventBook.Repository
	log   *zap.SugaredLogger
	errh  *evented.ErrLogger
}

func (server *EventQueryServer) GetEventBook(ctx context.Context, req *evented_query.Query) (*evented_core.EventBook, error) {
	id, err := evented_proto.ProtoToUUID(req.Root)
	if err != nil {
		return nil, err
	}
	var book evented_core.EventBook
	if req.LowerBound != 0 && req.UpperBound != 0 {
		book, err = server.repos.GetFromTo(ctx, id, req.LowerBound, req.UpperBound)
	} else if req.LowerBound != 0 {
		book, err = server.repos.GetFrom(ctx, id, req.LowerBound)
	} else {
		book, err = server.repos.Get(ctx, id)
	}
	return &book, nil
}

func (server *EventQueryServer) Listen(port uint16) {
	lis := support.CreateListener(port, server.errh)
	grpcServer := grpc.NewServer()

	evented_query.RegisterEventQueryServer(grpcServer, server)
	err := grpcServer.Serve(lis)
	server.errh.LogIfErr(err, "Failed starting server")
}
