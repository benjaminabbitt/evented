package eventQueryServer

import (
	eventedproto "github.com/benjaminabbitt/evented/proto"
	"github.com/benjaminabbitt/evented/proto/gen/github.com/benjaminabbitt/evented/proto/evented/core"
	"github.com/benjaminabbitt/evented/proto/gen/github.com/benjaminabbitt/evented/proto/evented/query"
	"github.com/benjaminabbitt/evented/repository/events"
	"github.com/benjaminabbitt/evented/support"
	"github.com/benjaminabbitt/evented/support/grpcWithInterceptors"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
)

func NewEventQueryServer(maxSize uint, repos events.EventStorer, log *zap.SugaredLogger) DefaultEventQueryServer {
	return DefaultEventQueryServer{
		EventBookSize: maxSize,
		eventRepos:    repos,
		log:           log,
	}
}

type DefaultEventQueryServer struct {
	query.UnimplementedEventQueryServer
	EventBookSize uint
	eventRepos    events.EventStorer
	log           *zap.SugaredLogger
}

func (o *DefaultEventQueryServer) GetEvents(req *query.Query, server query.EventQuery_GetEventsServer) error {
	id, err := eventedproto.ProtoToUUID(req.Root)
	if err != nil {
		return err
	}
	evtChan := make(chan *core.EventPage)
	var eventPages []*core.EventPage
	cover := &core.Cover{
		Domain: req.Domain,
		Root:   req.Root,
	}
	if req.LowerBound != 0 && req.UpperBound != 0 {
		err = o.eventRepos.GetFromTo(server.Context(), evtChan, id, req.LowerBound, req.UpperBound)
	} else if req.LowerBound != 0 {
		err = o.eventRepos.GetFrom(server.Context(), evtChan, id, req.LowerBound)
	} else {
		err = o.eventRepos.Get(server.Context(), evtChan, id)
	}
	maxSize := o.EventBookSize
	for page := range evtChan {
		pSize := uint(proto.Size(page))
		size := uint(0)
		// This approximation of size is not 100% correct, as of 20200415, it'll be about 2 bytes small per tests.
		// This addition is a performance optimization to avoid having to re-generate and re-serialize the event book repeatedly,
		//   and a single-digit-byte-class error isn't worth spending cycles on.
		if ((size + pSize) > maxSize) && (len(eventPages) > 0) {
			err := o.send(cover, eventPages, server)
			if err != nil {
				return err
			}
			size = 0
			eventPages = []*core.EventPage{}
		} else {
			size += pSize
			eventPages = append(eventPages, page)
		}
	}
	err = o.send(cover, eventPages, server)
	return nil
}

func (o *DefaultEventQueryServer) send(cover *core.Cover, pages []*core.EventPage, server query.EventQuery_GetEventsServer) error {
	book := &core.EventBook{
		Cover:    cover,
		Pages:    pages,
		Snapshot: nil,
	}
	err := server.Send(book)
	if err != nil {
		o.log.Error(err)
		return err
	}
	return nil
}

func (o *DefaultEventQueryServer) Synchronize(server query.EventQuery_SynchronizeServer) error {
	panic("implement me")
}

func (o *DefaultEventQueryServer) GetAggregateRoots(e *empty.Empty, server query.EventQuery_GetAggregateRootsServer) error {
	panic("implement me")
}

func (o *DefaultEventQueryServer) Listen(port uint, tracer opentracing.Tracer) error {
	lis := support.CreateListener(port, o.log)
	grpcServer := grpcWithInterceptors.GenerateConfiguredServer(o.log.Desugar(), tracer)

	query.RegisterEventQueryServer(grpcServer, o)
	err := grpcServer.Serve(lis)
	if err != nil {
		o.log.Error(err)
	}
	return err
}
