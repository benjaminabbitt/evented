package main

import (
	"context"
	evented_proto "github.com/benjaminabbitt/evented/proto"
	"github.com/benjaminabbitt/evented/proto/core"
	evented_query "github.com/benjaminabbitt/evented/proto/query"
	"github.com/benjaminabbitt/evented/repository/eventBook"
	"github.com/benjaminabbitt/evented/support"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func NewEventQueryServer(repos eventBook.EventBookStorer, log *zap.SugaredLogger) EventQueryServer {
	return EventQueryServer{
		repos: repos,
		log:   log,
	}
}

type EventQueryServer struct {
	evented_query.UnimplementedEventQueryServer
	repos eventBook.EventBookStorer
	log   *zap.SugaredLogger
}

func (o *EventQueryServer) GetEventBook(ctx context.Context, req *evented_query.Query) (*evented_core.EventBook, error) {
	id, err := evented_proto.ProtoToUUID(req.Root)
	if err != nil {
		return nil, err
	}
	var book *evented_core.EventBook
	if req.LowerBound != 0 && req.UpperBound != 0 {
		book, err = o.repos.GetFromTo(ctx, id, req.LowerBound, req.UpperBound)
	} else if req.LowerBound != 0 {
		book, err = o.repos.GetFrom(ctx, id, req.LowerBound)
	} else {
		book, err = o.repos.Get(ctx, id)
	}
	return book, nil
}

func (o *EventQueryServer) Listen(port uint16) error {
	lis := support.CreateListener(port, o.log)
	grpcServer := grpc.NewServer()

	evented_query.RegisterEventQueryServer(grpcServer, o)
	err := grpcServer.Serve(lis)
	if err != nil {
		o.log.Error(err)
	}
	return err
}
