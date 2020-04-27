package eventQueryServer

import (
	evented_proto "github.com/benjaminabbitt/evented/proto"
	evented_core "github.com/benjaminabbitt/evented/proto/core"
	evented_query "github.com/benjaminabbitt/evented/proto/query"
	"github.com/benjaminabbitt/evented/repository/events"
	"github.com/benjaminabbitt/evented/support"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/empty"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func NewEventQueryServer(maxSize uint, repos events.EventStorer, log *zap.SugaredLogger) DefaultEventQueryServer {
	return DefaultEventQueryServer{
		EventBookSize: maxSize,
		eventRepos:    repos,
		log:           log,
	}
}

type DefaultEventQueryServer struct {
	evented_query.UnimplementedEventQueryServer
	EventBookSize uint
	eventRepos    events.EventStorer
	log           *zap.SugaredLogger
}

func (o *DefaultEventQueryServer) GetEvents(req *evented_query.Query, server evented_query.EventQuery_GetEventsServer) error {
	id, err := evented_proto.ProtoToUUID(req.Root)
	if err != nil {
		return err
	}
	evtChan := make(chan *evented_core.EventPage)
	var eventPages []*evented_core.EventPage
	cover := &evented_core.Cover{
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
			eventPages = []*evented_core.EventPage{}
		} else {
			size += pSize
			eventPages = append(eventPages, page)
		}
	}
	err = o.send(cover, eventPages, server)
	return nil
}

func (o *DefaultEventQueryServer) send(cover *evented_core.Cover, pages []*evented_core.EventPage, server evented_query.EventQuery_GetEventsServer) error {
	book := &evented_core.EventBook{
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

func (o *DefaultEventQueryServer) Synchronize(server evented_query.EventQuery_SynchronizeServer) error {
	panic("implement me")
}

func (o *DefaultEventQueryServer) GetAggregateRoots(e *empty.Empty, server evented_query.EventQuery_GetAggregateRootsServer) error {
	panic("implement me")
}

func (o *DefaultEventQueryServer) Listen(port uint) error {
	lis := support.CreateListener(port, o.log)
	grpcServer := grpc.NewServer()

	evented_query.RegisterEventQueryServer(grpcServer, o)
	err := grpcServer.Serve(lis)
	if err != nil {
		o.log.Error(err)
	}
	return err
}
